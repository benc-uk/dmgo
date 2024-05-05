package gameboy

import (
	"log"
)

const ROM_BANK = 0x4000
const VRAM = 0x8000
const VRAM_CONT = 0x9000
const EXT_RAM = 0xA000
const WRAM = 0xC000
const ECHO_RAM = 0xE000
const OAM = 0xFE00
const OAM_END = 0xFE9F
const IO = 0xFF00
const HRAM = 0xFF80
const INT_ENABLE = 0xFFFF
const INT_FLAG = 0xFF0F

// HRAM register addresses
const LCD_CONTROL = 0xFF40
const LCD_STAT = 0xFF41
const SCROLL_Y = 0xFF42
const SCROLL_X = 0xFF43
const LCD_Y = 0xFF44
const BOOT_ROM_DISABLE = uint16(0xFF50)

const TILE_DATA_0 = 0x8000
const TILE_DATA_1 = 0x8800
const TILE_DATA_2 = 0x9000

const TILE_MAP_0 = 0x9800
const TILE_MAP_1 = 0x9C00

// Gameboy Memory Map
// 0x0000-0x3FFF: 16KB ROM Bank 00 (in cartridge, fixed at bank 00)
// 0x4000-0x7FFF: 16KB ROM Bank 01..NN (in cartridge, switchable bank number)
// 0x8000-0x9FFF: 8KB Video RAM (VRAM)
// 0xA000-0xBFFF: 8KB External RAM (in cartridge, switchable bank, if any)
// 0xC000-0xDFFF: 8KB Work RAM
// 0xE000-0xFDFF: 7.5KB Echo RAM - Reserved, Do Not Use
// 0xFE00-0xFE9F: 160B Sprite Attribute Table (OAM)
// 0xFEA0-0xFEFF: Not Usable
// 0xFF00-0xFF7F: 128B I/O Registers
// 0xFF80-0xFFFE: 127B High RAM (HRAM)
// 0xFFFF: Interrupt enable register

// Mapper is the memory map for the Gameboy
type Mapper struct {
	rom0      []byte
	rom1      []byte
	vram      []byte
	extRAM    []byte
	wram      []byte
	oam       []byte
	io        []byte
	hram      []byte
	interrupt byte

	bootROM       []byte
	bootROMLoaded bool

	ppu *PPU

	watches []uint16
}

func NewMapper() *Mapper {
	m := &Mapper{
		rom0:   make([]byte, 0x4000), // 16KB of ROM
		rom1:   make([]byte, 0x4000), // 16KB of ROM
		vram:   make([]byte, 0x2000), // 8KB of VRAM
		extRAM: make([]byte, 0x2000), // 8KB of external RAM
		wram:   make([]byte, 0x2000), // 8KB of WRAM
		oam:    make([]byte, 0x100),  // 160 bytes of OAM
		io:     make([]byte, 0x80),   // 128 bytes of IO
		hram:   make([]byte, 0x7F),   // 127 bytes of HRAM

		watches: []uint16{},

		bootROM:       make([]byte, 0x100),
		bootROMLoaded: false,
	}

	// Fill rom0 with 0xFF, really to simulate having no cartridge inserted
	for i := 0; i < 0x4000; i++ {
		m.rom0[i] = 0xFF
	}

	return m
}

func (m *Mapper) Write(addr uint16, data byte) {
	switch {
	case addr < ROM_BANK:
		{
			m.rom0[addr] = data
		}
	case addr >= ROM_BANK && addr < VRAM:
		{
			m.rom1[addr-ROM_BANK] = data
		}

	case addr >= VRAM && addr < EXT_RAM:
		{
			m.vram[addr-VRAM] = data

			// Check for writes to the tile data
			if addr >= TILE_DATA_0 && addr < TILE_MAP_0 && addr%16 == 15 {
				// Only update the tile if it's the 16th byte in the tile
				// TODO: This is more efficient but could it cause issues?
				m.ppu.updateTileCache(addr)
			}
		}

	case addr >= EXT_RAM && addr < WRAM:
		{
			m.extRAM[addr-EXT_RAM] = data
		}

	case addr >= WRAM && addr < ECHO_RAM:
		{
			m.wram[addr-WRAM] = data
		}

	case addr >= ECHO_RAM && addr < OAM:
		{
			// Spooky data written to the echo RAM is also written to the WRAM
			m.wram[addr-ECHO_RAM] = data
		}

	case addr >= OAM && addr <= OAM_END:
		{
			m.oam[addr-OAM] = data
			m.ppu.updateSpriteCache(addr)
		}

	case addr >= IO && addr < HRAM:
		{
			m.io[addr-IO] = data
		}

	case addr >= HRAM && addr < INT_ENABLE:
		{
			m.hram[addr-HRAM] = data
			// Special case for disabling the boot ROM
			if addr == BOOT_ROM_DISABLE && data == 0x01 {
				m.bootROMLoaded = false
			}
		}

	case addr == INT_ENABLE:
		{
			m.interrupt = data
		}

	case addr >= 0xFEA0 && addr <= 0xFEFF:
		{
			//log.Println("Invalid write to 0xFEA0-0xFEFF")
		}

	default:
		{
			log.Fatalf("Invalid memory write at %04X", addr)
		}
	}
}

func (m Mapper) Read(addr uint16) byte {
	switch {
	case addr < ROM_BANK:
		// Special case for the boot ROM overlay
		if m.bootROMLoaded && addr < 0x100 {
			return m.bootROM[addr]
		}

		return m.rom0[addr]
	case addr >= ROM_BANK && addr < VRAM:
		return m.rom1[addr-ROM_BANK]
	case addr >= VRAM && addr < EXT_RAM:
		return m.vram[addr-VRAM]
	case addr >= EXT_RAM && addr < WRAM:
		return m.extRAM[addr-EXT_RAM]
	case addr >= WRAM && addr < ECHO_RAM:
		return m.wram[addr-WRAM]
	case addr >= ECHO_RAM && addr < OAM:
		return m.wram[addr-ECHO_RAM]
	case addr >= OAM && addr <= OAM_END:
		return m.oam[addr-OAM]
	case addr >= IO && addr < HRAM:
		return m.io[addr-IO]
	case addr >= HRAM && addr < INT_ENABLE:
		return m.hram[addr-HRAM]
	case addr == INT_ENABLE:
		return m.interrupt
	// Invalid memory read
	case addr >= 0xFEA0 && addr <= 0xFEFF:
		{
			return 0xFF
		}
	}

	log.Fatalf("Invalid memory read at %04X", addr)

	return 0
}

func (m *Mapper) requestInterrupt(interrupt byte) {
	//log.Printf("Requesting interrupt %08b IE:%08b", interrupt, m.Read(INT_ENABLE))
	// Get the interrupt enable bit for the interrupt
	ie := m.Read(INT_ENABLE)
	if ie&interrupt == interrupt {
		// Set the interrupt flag
		m.Write(INT_FLAG, m.Read(INT_FLAG)|interrupt)
		//log.Printf("IF:%08b", m.Read(INT_FLAG))
	}
}
