package gameboy

import (
	"fmt"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	logging = false
)

const (
	INT_VBLANK = 0x01
	INT_LCD    = 0x02
	INT_TIMER  = 0x04
	INT_SERIAL = 0x08
	INT_JOYPAD = 0x10
)

type Gameboy struct {
	mapper *Mapper
	ppu    *PPU
	cpu    *CPU

	running bool
}

func NewGameboy() *Gameboy {
	mapper := NewMapper()
	cpu := NewCPU(mapper)

	gb := Gameboy{
		mapper:  mapper,
		ppu:     NewPPU(mapper),
		cpu:     cpu,
		running: false,
	}

	// Set up the initial state of the Gameboy
	mapper.Write(0xff00, 0x30) // Set the joypad register
	mapper.Write(0xff50, 0x00) // ENABLE the boot ROM
	mapper.Write(0xff40, 0x91) // Set the LCDC register
	mapper.Write(0xff41, 0x81) // Set the STAT register
	mapper.Write(0xff44, 0x90) // Set the scanline to 144
	mapper.Write(0xff47, 0xFC) // Set the background palette

	// gb.cpu.breakpoint = 0x028f
	gb.cpu.breakpoint = 0x0000
	gb.mapper.watches = []uint16{0xff02, 0xff80, 0xffa6}

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
	}

	if !mapper.bootROMLoaded {
		log.Println("Boot ROM not present, skipping to load cartridge ROM")
		mapper.Write(0xff50, 0x01) // DISABLE the boot ROM
		cpu.PC = 0x100
	}

	return &gb
}

// Update runs the system
func (gb *Gameboy) Update(cyclesPerFrame int) {
	if !gb.running {
		return
	}

	cycles := 0
	for cycles <= cyclesPerFrame {
		// Run the CPU fetch/exec cycle
		cpuCycles := gb.cpu.ExecuteNext()
		if cpuCycles <= 0 {
			log.Println("CPU halted")
			gb.DumpVRAM()
			gb.ppu.render()
			gb.Stop()
			break
		}

		cycles += cpuCycles

		// Handle interrupts
		if gb.cpu.IME {
			interrupt := gb.mapper.Read(INT_FLAG)
			enabled := gb.mapper.Read(INT_ENABLE)

			if interrupt&enabled > 0 {
				gb.cpu.handleInterrupt(interrupt & enabled)
			}
		}

		// PPU update
		gb.ppu.cycle(cpuCycles)
	}

	gb.ppu.render()

	// Read serial input
	hasData := gb.mapper.Read(0xff02)
	if hasData == 0x81 {
		// Read the data from the serial port
		data := gb.mapper.Read(0xff01)
		log.Printf("Serial data read: %c\n", data)
		gb.mapper.Write(0xff02, 0x00)
	}
}

func (gb *Gameboy) LoadROM(fileName string) {
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
	log.Printf("Read %d bytes from %s into ROM0\n", byteCount, fileName)

	if fileInfo.Size() > 0x4000 {
		// Read the next 16KB of the ROM into the second 16KB of the ROM space
		byteCount, err = file.Read(gb.mapper.rom1)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Read %d bytes from %s into ROM1\n", byteCount, fileName)
	}
}

func (gb *Gameboy) GetScreen() *ebiten.Image {
	return gb.ppu.screen
}

func (gb *Gameboy) Start() {
	gb.running = true
}

func (gb *Gameboy) Stop() {
	gb.running = false
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
	out += fmt.Sprintf("  LY: 0x%02X\n", gb.mapper.Read(0xff44))

	for _, addr := range gb.mapper.watches {
		out += fmt.Sprintf("Watch %04X:%02X\n", addr, gb.mapper.Read(addr))
	}

	return out
}

func SetLogging(l bool) {
	logging = l
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
