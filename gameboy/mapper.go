package gameboy

import (
	"log"
)

const ROM_BANK = 0x4000
const VRAM = 0x8000
const EXT_RAM = 0xA000
const WRAM = 0xC000
const ECHO_RAM = 0xE000
const OAM = 0xFE00
const NOT_USABLE = 0xFEA0
const IO = 0xFF00
const HRAM = 0xFF80

// IO Registers
const JOYP = 0xFF00
const SB = 0xFF01
const SC = 0xFF02
const DIV = 0xFF04
const TIMA = 0xFF05
const TMA = 0xFF06
const TAC = 0xFF07
const IF = 0xFF0F
const LCDC = 0xFF40
const STAT = 0xFF41
const SCY = 0xFF42
const SCX = 0xFF43
const LY = 0xFF44
const LYC = 0xFF45
const DMA = 0xFF46
const BGP = 0xFF47
const OBP0 = 0xFF48
const OBP1 = 0xFF49
const WY = 0xFF4A
const WX = 0xFF4B
const BOOT_ROM_DISABLE = 0xFF50
const IE = 0xFFFF

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

	bootROM []byte

	watches []uint16
	buttons *Buttons
}

func NewMapper(buttons *Buttons) *Mapper {
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
		buttons: buttons,
	}

	// Fill rom0 with 0xFF, really just to simulate having no cartridge inserted
	for i := 0; i < 0x4000; i++ {
		m.rom0[i] = 0xFF
	}

	return m
}

func (m *Mapper) write(addr uint16, data byte) {
	switch {
	case addr < ROM_BANK:
		{
			// Ignore writes to the ROM bank 0
			//m.rom0[addr] = data
		}
	case addr >= ROM_BANK && addr < VRAM:
		{
			// Ignore writes to the ROM bank 1
			//m.rom1[addr-ROM_BANK] = data
		}

	case addr >= VRAM && addr < EXT_RAM:
		{
			m.vram[addr-VRAM] = data

			// Check for writes to the tile data
			// if addr >= TILE_DATA_0 && addr < TILE_MAP_0 {
			// 	m.ppu.updateTileCache(addr)
			// }
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

	case addr >= OAM && addr < NOT_USABLE:
		{
			m.oam[addr-OAM] = data
			//m.ppu.updateSpriteCache(addr)
		}

	case addr >= NOT_USABLE && addr < IO:
		{
			// Writing to this range is prohibited, but code often tries
			return
		}

	case addr >= IO && addr < HRAM:
		{
			if addr == DIV {
				// Writing to the DIV register resets the counter
				m.io[addr-IO] = 0
				return
			}

			if addr == DMA {
				m.io[addr-IO] = data
				for i := 0; i < 0xA0; i++ {
					// For DMA source address is divided by 0x100 for some reason
					m.write(OAM+uint16(i), m.read(uint16(data)*0x100+uint16(i)))
				}
				return
			}

			// if addr == SC {
			// 	if data == 0x81 {
			// 		log.Printf("Serial transfer: %c\n", m.io[01])

			// 		// Reset the transfer flag
			// 		m.io[addr-IO] = 0x80
			// 	}
			// 	return
			// }

			m.io[addr-IO] = data
		}

	case addr >= HRAM && addr < IE:
		{
			m.hram[addr-HRAM] = data
		}

	case addr == IE:
		{
			m.interrupt = data
		}

	default:
		{
			log.Fatalf("Invalid memory write at %04X", addr)
		}
	}
}

func (m Mapper) read(addr uint16) byte {
	switch {
	case addr < ROM_BANK:
		// Special case for the boot ROM, which is overlaid on the first 256 bytes of memory
		if addr < 0x100 && m.bootROMEnabled() {
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
	case addr >= OAM && addr < NOT_USABLE:
		return m.oam[addr-OAM]
		// Reading from the NOT_USABLE range returns 0xFF
	case addr >= NOT_USABLE && addr < IO:
		return 0xFF
	case addr >= IO && addr < HRAM:
		// Handle requests for the JOYP register
		if addr == JOYP {
			// If bit 4 is NOT set, return the d-pad state
			if m.io[0]&0x10 == 0 {
				return m.io[0] | m.buttons.getPadState()
			}

			// If bit 5 is ZERO, return the button state
			if m.io[0]&0x20 == 0 {
				return m.io[0] | m.buttons.getButtonState()
			}

			// If neither bit 4 or 5 are set, return 0xFF
			return m.io[0] | 0x0F
		}

		// if addr == STAT {
		// 	log.Printf("Reading STAT register\n")
		// }

		// if addr == LCDC {
		// 	log.Printf("Reading LCDC register\n")
		// }

		// if addr == LY {
		// 	//log.Printf("Reading LY register %02X\n", m.io[68])
		// 	return 153
		// }

		// Rest of the IO registers are just read
		return m.io[addr-IO]
	case addr >= HRAM && addr < IE:
		return m.hram[addr-HRAM]
	case addr == IE:
		return m.interrupt

	}

	log.Fatalf("Invalid memory read at %04X", addr)

	return 0
}

func (m *Mapper) bootROMEnabled() bool {
	return len(m.bootROM) > 0 && m.read(BOOT_ROM_DISABLE) == 0
}

func (m *Mapper) loadBootROM(data []byte) {
	log.Printf("Configuring boot ROM")
	if len(data) != 0x100 {
		log.Fatalf("Boot ROM is not the correct size, got %d bytes, expected 256 bytes", len(data))
	}

	m.write(BOOT_ROM_DISABLE, 0x00) // ENABLE the boot ROM
	m.bootROM = data
}
