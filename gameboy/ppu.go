package gameboy

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// VRAM offests
const TILE_DATA_0 = 0x8000
const TILE_DATA_1 = 0x8800
const TILE_DATA_2 = 0x9000
const TILE_MAP_0 = uint16(0x9800)
const TILE_MAP_1 = uint16(0x9C00)

type PPU struct {
	mapper     *Mapper
	emuPalette [4]color.RGBA
	screen     *ebiten.Image
	tileCache  map[uint16]*ebiten.Image

	// Scanline register
	scanline   byte
	dotCounter int

	gb *Gameboy
}

// Abstraction over the OAM data to represent a sprite
type Sprite struct {
	y        byte
	x        byte
	tile     byte
	priority bool
	flipY    bool
	flipX    bool
}

func NewPPU(mapper *Mapper) *PPU {
	// Hard coded palette for now
	pallet := [4]color.RGBA{
		{0xe0, 0xf8, 0xd0, 255},
		{0x88, 0xc0, 0x70, 255},
		{0x34, 0x68, 0x56, 255},
		{0x08, 0x18, 0x20, 255},
	}

	ppu := &PPU{
		emuPalette: pallet,
		mapper:     mapper,

		// Internal screen buffer
		screen: ebiten.NewImage(160, 144),
	}

	return ppu
}

func (ppu *PPU) newSprite(addr uint16) Sprite {
	i := (addr - OAM) / 4

	sprite := Sprite{
		y:        ppu.mapper.oam[i*4],
		x:        ppu.mapper.oam[i*4+1],
		tile:     ppu.mapper.oam[i*4+2],
		priority: ppu.mapper.oam[i*4+3]&0x80 == 0x80,
		flipY:    ppu.mapper.oam[i*4+3]&0x40 == 0x40,
		flipX:    ppu.mapper.oam[i*4+3]&0x20 == 0x20,
	}

	return sprite
}

// This function creates a new 8x8 image for a tile reading 16 bytes from the VRAM
// and using the BGP register to lookup the palette
func (ppu *PPU) getTileImage(addr uint16) *ebiten.Image {
	// Check if the tile is already in the cache
	if _, ok := ppu.tileCache[addr]; ok {
		return ppu.tileCache[addr]
	}

	pixels := make([]byte, 8*8*4)

	// Read the Gameboy palette state, this can be changed by the game
	palLookup := ppu.mapper.read(BGP)

	for tileByteIndex := uint16(0); tileByteIndex < 16; tileByteIndex += 2 {
		byte1 := ppu.mapper.read(addr + tileByteIndex)
		byte2 := ppu.mapper.read(addr + tileByteIndex + 1)
		y := int(tileByteIndex / 2)
		for bit := 0; bit < 8; bit++ {
			// Combine the bits to get the color index
			colorId := (byte1 >> (7 - bit) & 1) | ((byte2 >> (7 - bit) & 1) << 1)

			// Use the palette BGP, which is byte with 2bit colorId -> Value mapping
			// https://gbdev.io/pandocs/Palettes.html
			colorVal := palLookup >> (colorId * 2) & 0x3

			// Set the final color in the pixel array
			pixels[(y*8+bit)*4] = ppu.emuPalette[colorVal].R
			pixels[(y*8+bit)*4+1] = ppu.emuPalette[colorVal].G
			pixels[(y*8+bit)*4+2] = ppu.emuPalette[colorVal].B
			pixels[(y*8+bit)*4+3] = 255
		}
	}

	img := ebiten.NewImage(8, 8)
	img.WritePixels(pixels)

	// Cache the tile
	ppu.tileCache[addr] = img
	return img
}

func (ppu *PPU) getTileAddr(tileNum byte) uint16 {
	// Addressing mode 8000 is sane and normal
	if ppu.GetLCDCBit(4) == 1 {
		return uint16(TILE_DATA_0 + uint16(tileNum)*16)
	}

	// Tile number is a *signed* 8bit when in 8800 address mode
	return uint16(int(TILE_DATA_2) + int(int8(tileNum))*16)
}

func (ppu *PPU) render() {
	mapBase := TILE_MAP_0
	if ppu.GetLCDCBit(3) == 1 {
		mapBase = TILE_MAP_1
	}

	// Reset cache on each render, coz palette can change
	ppu.tileCache = make(map[uint16]*ebiten.Image)

	// get SCROLL_Y and SCROLL_X
	scrollY := float64(ppu.mapper.read(SCY))
	scrollX := float64(ppu.mapper.read(SCX))

	// Read the 1024 bytes of tile map data
	// And render into the screen at the correct position
	for i := uint16(0); i < 1024; i++ {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64((i%32)*8)+scrollX, float64((i/32)*8)-scrollY)

		tilenum := int(ppu.mapper.read(mapBase + i))
		tileAddr := ppu.getTileAddr(byte(tilenum))
		ppu.screen.DrawImage(ppu.getTileImage(tileAddr), op)
	}

	// Handle OAM and render 40 sprites
	for i := 0; i < 40; i++ {
		addr := OAM + uint16(i*4)
		sprite := ppu.newSprite(addr)
		if sprite.y == 0 && sprite.x == 0 {
			continue
		}

		op := &ebiten.DrawImageOptions{}
		screenY := int(sprite.y) - 16
		screenX := int(sprite.x) - 8
		op.GeoM.Translate(float64(screenX), float64(screenY))

		// Tile addressing is more simple for sprites
		tileAddr := TILE_DATA_0 + uint16(sprite.tile)*16

		// Flip the sprite if needed
		if sprite.flipX {
			op.GeoM.Scale(-1, 1)
			op.GeoM.Translate(8, 0)
		}

		if sprite.flipY {
			op.GeoM.Scale(1, -1)
			op.GeoM.Translate(0, 8)
		}

		ppu.screen.DrawImage(ppu.getTileImage(tileAddr), op)
	}
}

func (ppu *PPU) cycle(clockCycles int) {
	ppu.dotCounter += clockCycles

	// TODO: Still needs work
	if ppu.dotCounter > 456 {
		ppu.dotCounter = 0

		ppu.scanline++

		if ppu.scanline == 144 {
			// Request vblank interrupt
			ppu.gb.requestInterrupt(INT_VBLANK)
		}

		if ppu.scanline > 153 {
			ppu.scanline = 0
			// Render has been moved to the main GB loop
		}

		ppu.mapper.write(LY, ppu.scanline)
		if ppu.scanline == ppu.mapper.read(LYC) {
			// set bit 2 of STAT
			ppu.mapper.write(STAT, bitSet(ppu.mapper.read(STAT), 2))
		}
	}
}

// LCD Control Register
// Bit 7 - LCD Display Enable (0=Off, 1=On)
// Bit 6 - Window Tile Map Display Select (0=9800-9BFF, 1=9C00-9FFF)
// Bit 5 - Window Display Enable (0=Off, 1=On)
// Bit 4 - Tile select for background (0=9000-97FF, 1=8000-8FFF)
// Bit 3 - Tile Data Select (0=9800-9BFF, 1=9C00-9FFF)
// Bit 2 - OBJ Size (0=8x8, 1=8x16)
// Bit 1 - OBJ Display Enable (0=Off, 1=On)
// Bit 0 - BG/Window Display/Priority (0=Off, 1=On)
func (ppu *PPU) GetLCDCBit(bit byte) byte {
	return ppu.mapper.read(LCDC) >> bit & 1
}
