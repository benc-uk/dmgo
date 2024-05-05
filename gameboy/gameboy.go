package gameboy

import (
	"fmt"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

var logging = false

const (
	INT_VBLANK = 0x01
	INT_LCD    = 0x02
	INT_TIMER  = 0x04
	INT_SERIAL = 0x08
	INT_JOYPAD = 0x10
)

type Config struct {
	Scale       int      `yaml:"scale"`
	ROM         string   `yaml:"rom"`
	Logging     bool     `yaml:"logging"`
	Breakpoints []uint16 `yaml:"breakpoints"`
	Watches     []uint16 `yaml:"watches"`
	BootROM     bool     `yaml:"bootROM"`
}

type Gameboy struct {
	mapper *Mapper
	ppu    *PPU
	cpu    *CPU

	divider int

	Running bool
	config  Config
}

func NewGameboy(config Config) *Gameboy {
	mapper := NewMapper()
	cpu := NewCPU(mapper)

	gb := Gameboy{
		mapper:  mapper,
		ppu:     NewPPU(mapper),
		cpu:     cpu,
		Running: false,

		config: config,
	}

	// Set up the initial state of the Gameboy
	mapper.io[0x00] = 0xf
	mapper.Write(0xff50, 0x00) // ENABLE the boot ROM
	mapper.Write(0xff40, 0x91) // Set the LCDC register
	mapper.Write(0xff41, 0x81) // Set the STAT register
	mapper.Write(0xff44, 0x90) // Set the scanline to 144
	mapper.Write(0xff47, 0xE4) // Set the background palette

	if config.BootROM {
		// Load the boot ROM from res/dmg_boot.bin
		// If it fails, the boot ROM will be skipped and not used
		mapper.bootROMLoaded = true
		bootROMFile, err := os.Open("res/dmg_boot.bin")
		if err != nil {
			mapper.bootROMLoaded = false
		} else {
			_, err = bootROMFile.Read(mapper.bootROM)
			if err != nil {
				mapper.bootROMLoaded = false
			}
			log.Println("Boot ROM loaded OK")
		}
	}

	if !mapper.bootROMLoaded {
		log.Println("Boot ROM not loaded, it will be disabled")
		mapper.Write(0xff50, 0x01) // DISABLE the boot ROM
		// Set the initial PC to 0x100 which is where the game ROM starts
		cpu.PC = 0x100
	}

	logging = config.Logging

	if len(config.Breakpoints) > 0 {
		cpu.breakpoints = config.Breakpoints
	}

	if len(config.Watches) > 0 {
		mapper.watches = config.Watches
	}

	return &gb
}

// Update runs the system
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
		// Run the CPU fetch/exec cycle
		cpuCycles := gb.cpu.ExecuteNext(false)
		if cpuCycles < 0 {
			log.Println("Stopping emulation due to error")
			gb.Running = false
			break
		}

		cycles += cpuCycles
		cycles += gb.checkInterrupts()

		// PPU update
		gb.ppu.cycle(cpuCycles)
	}

	// TODO: Remove this later I think
	gb.ppu.render()

	// Timer update DIV
	// TODO: This is not correct yet!
	gb.updateTimers(cycles)

	// Read serial input
	hasData := gb.mapper.Read(0xff02)
	if hasData == 0x81 {
		// Read the data from the serial port
		data := gb.mapper.Read(0xff01)
		log.Printf("Serial data read: %c\n", data)
		gb.mapper.Write(0xff02, 0x00)
	}
}

func (gb *Gameboy) checkInterrupts() int {
	// Check for interrupts
	if gb.cpu.IME {
		interrupts := gb.mapper.Read(INT_FLAG) & gb.mapper.Read(INT_ENABLE)
		if interrupts != 0 {
			if logging {
				log.Printf("Interrupts: %08b\n", interrupts)
			}

			if interrupts&INT_VBLANK != 0 {
				gb.cpu.handleInterrupt(INT_VBLANK)
				return 20
			}
			if interrupts&INT_LCD != 0 {
				gb.cpu.handleInterrupt(INT_LCD)
				return 20
			}
			if interrupts&INT_TIMER != 0 {
				gb.cpu.handleInterrupt(INT_TIMER)
				return 20
			}
			if interrupts&INT_SERIAL != 0 {
				gb.cpu.handleInterrupt(INT_SERIAL)
				return 20
			}
			if interrupts&INT_JOYPAD != 0 {
				gb.cpu.handleInterrupt(INT_JOYPAD)
				return 20
			}
		}
	}

	return 0
}

func (gb *Gameboy) updateTimers(cycles int) {
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
	out += fmt.Sprintf("PC: 0x%04X -> %s\n\n", gb.cpu.PC, opcodeNames[gb.mapper.Read(cpu.PC)])
	out += fmt.Sprintf("A:%02X B:%02X C:%02X D:%02X E:%02X H:%02X L:%02X\n",
		cpu.A(), cpu.B(), cpu.C(), cpu.D(), cpu.E(), cpu.H(), cpu.L())
	out += fmt.Sprintf("AF:%04X BC:%04X DE:%04X HL:%04X SP:%04X\n", cpu.AF, cpu.BC, cpu.DE, cpu.HL, cpu.SP)
	out += fmt.Sprintf("IE:%08b IF:%08b IME:%d\n", gb.mapper.Read(INT_ENABLE), gb.mapper.Read(INT_FLAG), BoolToInt(cpu.IME))

	// Flags
	out += fmt.Sprintf("Z:%d N:%d H:%d C:%d\n\n",
		BoolToInt(cpu.getFlagZ()), BoolToInt(cpu.getFlagN()), BoolToInt(cpu.getFlagH()), BoolToInt(cpu.getFlagC()))

	// Show the next 5 bytes of memory
	for i := cpu.PC; i < cpu.PC+5; i++ {
		out += fmt.Sprintf("%04X: 0x%02X\n", i, gb.mapper.Read(i))
	}

	out += fmt.Sprintf("\nLCDC: 0x%02X %08b\n", gb.mapper.Read(0xff40), gb.mapper.Read(0xff40))
	out += fmt.Sprintf("STAT: 0x%02X\n", gb.mapper.Read(0xff41))
	out += fmt.Sprintf("  LY: 0x%02X\n\n", gb.mapper.Read(0xff44))

	for _, addr := range gb.mapper.watches {
		out += fmt.Sprintf("Watch %04X:%02X\n", addr, gb.mapper.Read(addr))
	}

	return out
}

func (gb *Gameboy) DumpVRAM() {
	file, err := os.Create("vram.dump")
	if err != nil {
		log.Fatal(err)
	}

	for i := uint16(0x8000); i < 0x9FFF; i++ {
		file.Write([]byte{gb.mapper.Read(i)})
	}
}
