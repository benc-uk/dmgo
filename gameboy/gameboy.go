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
func (gb *Gameboy) Cycle() {
	if !gb.running {
		return
	}

	gb.dots++

	// PPU runs every two dots
	gb.ppu.cycle()

	// CPU runs every four dots
	if gb.dots%2 == 0 {
		// Run the CPU fetch/exec cycle
		ok := gb.cpu.ExecuteNext()
		if !ok {
			gb.ppu.render() // TODO: Remove this
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

	out := fmt.Sprintf("Prev: " + gb.cpu.opDebug)
	out += fmt.Sprintf("\nPC: 0x%04x\n\n", gb.cpu.PC)
	out += fmt.Sprintf("A: 0x%02X B: 0x%02X C: 0x%02X D: 0x%02X\nE: 0x%02X H: 0x%02X L: 0x%02X SP: 0x%04X\n",
		cpu.getRegA(), cpu.getRegB(), cpu.getRegC(), cpu.getRegD(), cpu.getRegE(), cpu.getRegH(), cpu.getRegL(), cpu.SP)

	// Flags
	out += fmt.Sprintf("Z: %d N: %d H: %d C: %d\n\n",
		BoolToInt(cpu.getFlagZ()), BoolToInt(cpu.getFlagN()), BoolToInt(cpu.getFlagH()), BoolToInt(cpu.getFlagC()))

	// Show the next 10 bytes of memory
	for i := cpu.PC; i < cpu.PC+10; i++ {
		out += fmt.Sprintf("%04X: 0x%02X\n", i, gb.mapper.Read(i))
	}

	out += fmt.Sprintf("\n%04X: 0x%02b\n", 0xff40, gb.mapper.Read(0xff40))
	out += fmt.Sprintf("%04X: 0x%02b\n", 0xff41, gb.mapper.Read(0xff41))
	out += fmt.Sprintf("%04X: 0x%02X\n", 0xff44, gb.mapper.Read(0xff44))
	out += fmt.Sprintf("%04X: 0x%02b\n", 0xff47, gb.mapper.Read(0xff47))

	return out
}

func SetLogging(l bool) {
	logging = l
}
