package gameboy

const ROM_BANK = 0x4000
const VRAM = 0x8000
const VRAM_CONT = 0x9000
const EXT_RAM = 0xA000
const OAM = 0xFE00
const OAM_END = 0xFE9F
const IO = 0xFF00
const HRAM = 0xFF80

// HRAM register addresses
const LCD_CONTROL = 0xFF40
const LCD_STAT = 0xFF41
const LCD_Y = 0xFF44

const TILE_DATA_1 = 0x8000
const TILE_DATA_2 = 0x8800
const TILE_DATA_3 = 0x9000

const TILE_MAP_1 = 0x9800
const TILE_MAP_2 = 0x9C00

type Mapper struct {
	rom0       []byte
	rom1       []byte
	vram       []byte
	extRAM     []byte
	wram       []byte
	oam        []byte
	io         []byte
	hram       []byte
	interrupts byte

	ppu *PPU
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
	}

	return m
}

func (m Mapper) Write(addr uint16, data byte) {
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
			if addr >= TILE_DATA_1 && addr < TILE_MAP_1 && addr%16 == 15 {
				// Only update the tile if it's the 16th byte in the tile
				// TODO: This is more efficient but could it cause issues?
				m.ppu.updateTileCache(addr)
			}
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
	}
}

func (m Mapper) Read(addr uint16) byte {
	switch {
	case addr < ROM_BANK:
		return m.rom0[addr]
	case addr >= ROM_BANK && addr < VRAM:
		return m.rom1[addr-ROM_BANK]
	case addr >= VRAM && addr < EXT_RAM:
		return m.vram[addr-VRAM]
	case addr >= OAM && addr <= OAM_END:
		return m.oam[addr-OAM]
	case addr >= IO && addr < HRAM:
		return m.io[addr-IO]
	}

	return 0
}
