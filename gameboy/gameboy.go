package gameboy

import (
	"fmt"
	"log"
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	logging = true
)

type Gameboy struct {
	mapper *Mapper
	ppu    *PPU
	cpu    *CPU

	running bool
	dots    int64
}

func NewGameboy() *Gameboy {
	mapper := NewMapper()

	gb := Gameboy{
		mapper:  mapper,
		ppu:     NewPPU(mapper),
		cpu:     NewCPU(mapper),
		running: false,
		dots:    0,
	}

	mapper.Write(0xff50, 0x01) // Disable the boot ROM
	mapper.Write(0xff40, 0x91) // Set the LCDC register
	mapper.Write(0xff41, 0x81) // Set the STAT register
	mapper.Write(0xff44, 0x90) // Set the scanline to 144
	mapper.Write(0xff47, 0xFC) // Set the background palette

	// gb.cpu.breakpoint = 0x078C

	return &gb
}

func (gb *Gameboy) LoadMemDump(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}

	memDump := make([]byte, 0xffff)
	byteCount, err := file.Read(memDump)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Read %d bytes from %s\n", byteCount, fileName)

	for i := uint16(0); i < 0xffff; i++ {
		gb.mapper.Write(i, memDump[i])
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

// Cycle runs the system for two dot cycles (2 ticks of 4.19MHz)
func (gb *Gameboy) Cycle(force bool) {
	if force {
		gb.ppu.cycle()
		_ = gb.cpu.ExecuteNext(true)
	}

	if !gb.running {
		return
	}

	gb.dots++

	// PPU runs every two dots
	gb.ppu.cycle()

	// CPU runs every four dots
	if gb.dots%2 == 0 {
		// Run the CPU fetch/exec cycle
		ok := gb.cpu.ExecuteNext(false)
		if !ok {
			gb.ppu.render() // TODO: Remove this
			gb.DumpVRAM()
			gb.Stop()
		}
	}

	// I have no idea if this is a real risk!
	if gb.dots > math.MaxInt64-1 {
		gb.dots = 0
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

	out := fmt.Sprintf("Last instr: " + gb.cpu.opDebug)
	out += fmt.Sprintf("\nPC: 0x%04X\n\n", gb.cpu.PC)
	out += fmt.Sprintf("A: %02X B: %02X C: %02X D: %02X\nE: %02X H: %02X L: %02X SP: %04X\n",
		cpu.A(), cpu.B(), cpu.C(), cpu.D(), cpu.E(), cpu.H(), cpu.L(), cpu.SP)
	out += fmt.Sprintf("AF: %04X BC: %04X DE: %04X HL: %04X\n\n", cpu.AF, cpu.BC, cpu.DE, cpu.HL)

	// Flags
	out += fmt.Sprintf("Z: %d N: %d H: %d C: %d\n\n",
		BoolToInt(cpu.getFlagZ()), BoolToInt(cpu.getFlagN()), BoolToInt(cpu.getFlagH()), BoolToInt(cpu.getFlagC()))

	// Show the next 5 bytes of memory
	for i := cpu.PC; i < cpu.PC+5; i++ {
		out += fmt.Sprintf("%04X: 0x%02X\n", i, gb.mapper.Read(i))
	}

	out += fmt.Sprintf("\n%04X: 0x%02b\n", 0xff40, gb.mapper.Read(0xff40))
	out += fmt.Sprintf("%04X: 0x%02b\n", 0xff41, gb.mapper.Read(0xff41))
	out += fmt.Sprintf("%04X: 0x%02X\n", 0xff44, gb.mapper.Read(0xff44))
	out += fmt.Sprintf("%04X: 0x%02b\n\n", 0xff47, gb.mapper.Read(0xff47))

	// Stack
	for i := cpu.SP - 2; i < cpu.SP+5; i++ {
		out += fmt.Sprintf("%04X: 0x%02X\n", i, gb.mapper.Read(i))
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

	for i := uint16(0x8000); i < 0x9800; i++ {
		file.Write([]byte{gb.mapper.Read(i)})
	}
}
