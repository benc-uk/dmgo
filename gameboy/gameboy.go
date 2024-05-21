package gameboy

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	INT_VBLANK = 0x01
	INT_LCD    = 0x02
	INT_TIMER  = 0x04
	INT_SERIAL = 0x08
	INT_JOYPAD = 0x10
)

type Config struct {
	BootROM     string   `yaml:"bootROM"`
	Breakpoints []uint16 `yaml:"breakpoints"`
	Watches     []uint16 `yaml:"watches"`
	OpcodeDebug []byte   `yaml:"opcodeDebug"`
}

type Gameboy struct {
	mapper  *Mapper
	ppu     *PPU
	cpu     *CPU
	Buttons *Buttons

	divider int

	Running bool
	config  Config
}

func NewGameboy(config Config) *Gameboy {
	log.Println("Initializing Gameboy")

	buttons := &Buttons{}
	mapper := NewMapper(buttons)
	cpu := NewCPU(mapper)
	ppu := NewPPU(mapper)
	gb := Gameboy{
		mapper:  mapper,
		ppu:     ppu,
		cpu:     cpu,
		Running: false,

		config:  config,
		Buttons: buttons,
	}

	ppu.gb = &gb // Ugly cross dependency, so PPU can request interrupts

	// Set up the initial state of the Gameboy
	mapper.write(LCDC, 0x91) // Set the LCDC register
	mapper.write(STAT, 0x81) // Set the STAT register
	mapper.write(LY, 0x91)   // Set the scanline to 145
	mapper.write(BGP, 0xFC)  // Set the background palette
	mapper.write(DIV, 0xAB)  // Set the divider

	// Optional boot ROM, not needed but included for authenticity
	if config.BootROM != "" {
		bootROMFile, err := os.Open(config.BootROM)
		if err != nil {
			log.Fatal(err)
		}

		br, err := io.ReadAll(bootROMFile)
		if err != nil {
			log.Fatal(err)
		}

		mapper.loadBootROM(br)
	}

	if !mapper.bootROMEnabled() {
		log.Println("Boot ROM not available, it will be disabled")
		// DISABLE the boot ROM
		mapper.write(BOOT_ROM_DISABLE, 0x01)
		// Jump PC to 0x100, this is where PC would be after the boot ROM
		cpu.pc = 0x100
	}

	if len(config.Breakpoints) > 0 {
		cpu.breakpoints = config.Breakpoints
	}

	if len(config.Watches) > 0 {
		mapper.watches = config.Watches
	}

	if len(config.OpcodeDebug) > 0 {
		cpu.opDebug = config.OpcodeDebug
	}

	return &gb
}

// Update runs the system each frame
func (gb *Gameboy) Update(cyclesPerFrame int) {
	// This is how we step manually
	if cyclesPerFrame <= 0 {
		cpuCycles := gb.cpu.ExecuteNext(true)
		cpuCycles += gb.checkInterrupts()

		// PPU update
		gb.ppu.cycle(cpuCycles)

		return
	}

	if !gb.Running {
		return
	}

	cycles := 0
	for cycles <= cyclesPerFrame {
		cycles += gb.checkInterrupts()

		// Run the CPU fetch/exec cycle
		cpuCycles := gb.cpu.ExecuteNext(false)
		if cpuCycles < 0 {
			log.Println("Stopping emulation due to error")
			gb.Running = false
			break
		}

		cycles += cpuCycles

		// PPU update
		gb.ppu.cycle(cpuCycles)

		// Read serial port, really only used for debugging and Blargg's tests
		if gb.mapper.io[0x02] == 0x81 {
			fmt.Printf("%c", gb.mapper.io[0x01])
			gb.mapper.io[0x02] = 0x80
		}
	}

	// Timer update DIV
	gb.updateTimers(cycles)

	// Interrupt for joypad
	if gb.Buttons.Changed() {
		gb.requestInterrupt(INT_JOYPAD)
		gb.Buttons.ClearChanged()
	}

	// HACK: Not sure this needs to be here
	gb.ppu.render()
}

func (gb *Gameboy) checkInterrupts() int {

	if gb.cpu.halted && !gb.cpu.ime {
		return 0
	}

	// Check for interrupts
	if gb.mapper.read(IF)&gb.mapper.read(IE) != 0 {
		gb.cpu.halted = false
	}

	if gb.cpu.ime {
		interruptMask := gb.mapper.read(IF) & gb.mapper.read(IE)

		if interruptMask != 0 {
			if gb.cpu.halted {
				log.Printf("Halted, and interrupt %08b requested, pc: 0x%04X", interruptMask, gb.cpu.pc)
				gb.cpu.halted = false
			}

			if interruptMask&INT_VBLANK != 0 {
				gb.cpu.handleInterrupt(INT_VBLANK)
				return 20
			}
			if interruptMask&INT_LCD != 0 {
				gb.cpu.handleInterrupt(INT_LCD)
				return 20
			}
			if interruptMask&INT_TIMER != 0 {
				gb.cpu.handleInterrupt(INT_TIMER)
				return 20
			}
			if interruptMask&INT_SERIAL != 0 {
				gb.cpu.handleInterrupt(INT_SERIAL)
				return 20
			}
			if interruptMask&INT_JOYPAD != 0 {
				gb.cpu.handleInterrupt(INT_JOYPAD)
				return 20
			}
		}
	}

	return 0
}

func (gb *Gameboy) requestInterrupt(interruptBit byte) {
	interruptByte := gb.mapper.read(IF)
	interruptByte |= interruptBit
	gb.mapper.write(IF, interruptByte)
}

func (gb *Gameboy) updateTimers(cycles int) {
	// TODO: Not sure this is the correct way to handle the DIV register
	gb.divider += cycles
	if gb.divider >= 256 {
		gb.divider -= 256
		// Note that the DIV register is actually at 0xff04
		// We can't use Write() as that resets the divider
		gb.mapper.io[DIV-IO]++
	}
}

func (gb *Gameboy) LoadROM(fileName string) {
	log.Printf("Loading ROM: %s\n", fileName)
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	fileInfo, _ := file.Stat()

	// Read the first 16KB of the ROM into the first 16KB of the ROM space
	byteCount, err := file.Read(gb.mapper.rom0)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Read %d bytes into ROM0\n", byteCount)

	if fileInfo.Size() > 0x4000 {
		// Read the next 16KB of the ROM into the second 16KB of the ROM space
		byteCount, err = file.Read(gb.mapper.rom1)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Read %d bytes into ROM1\n", byteCount)
	}
}

func (gb *Gameboy) GetScreen() *ebiten.Image {
	return gb.ppu.screen
}

func (gb *Gameboy) GetDebugInfo() string {
	cpu := gb.cpu

	out := ""
	out += fmt.Sprintf("PC: 0x%04X -> %s\n\n", gb.cpu.pc, opcodeNames[gb.mapper.read(cpu.pc)])
	out += fmt.Sprintf("A:%02X B:%02X C:%02X D:%02X E:%02X H:%02X L:%02X\n",
		cpu.A(), cpu.B(), cpu.C(), cpu.D(), cpu.E(), cpu.H(), cpu.L())
	out += fmt.Sprintf("AF:%04X BC:%04X DE:%04X HL:%04X SP:%04X\n", cpu.af, cpu.bc, cpu.de, cpu.hl, cpu.sp)
	out += fmt.Sprintf("IE:%08b IF:%08b IME:%d\n", gb.mapper.read(IE), gb.mapper.read(IF), BoolToInt(cpu.ime))

	// Flags
	out += fmt.Sprintf("Z:%d N:%d H:%d C:%d\n\n",
		BoolToInt(cpu.getFlagZ()), BoolToInt(cpu.getFlagN()), BoolToInt(cpu.getFlagH()), BoolToInt(cpu.getFlagC()))

	// Show the next 4 bytes of memory
	for i := cpu.pc; i < cpu.pc+4; i++ {
		out += fmt.Sprintf("%04X: 0x%02X\n", i, gb.mapper.read(i))
	}

	out += fmt.Sprintf("\nLCDC: 0x%08b\n", gb.mapper.read(LCDC))
	out += fmt.Sprintf("STAT: %08b\n", gb.mapper.read(STAT))
	out += fmt.Sprintf("  LY: 0x%02X\n", gb.mapper.read(LY))
	out += fmt.Sprintf(" LYC: 0x%02X\n", gb.mapper.read(LYC))
	out += fmt.Sprintf(" BGP: %08b\n", gb.mapper.read(BGP))
	out += fmt.Sprintf("  SB: %08b 0x%02X\n", gb.mapper.io[1], gb.mapper.io[1])
	out += fmt.Sprintf("  SC: %08b 0x%02X\n\n", gb.mapper.io[2], gb.mapper.io[2])

	for _, addr := range gb.mapper.watches {
		out += fmt.Sprintf("Watch %04X:%02X\n", addr, gb.mapper.read(addr))
	}

	return out
}
