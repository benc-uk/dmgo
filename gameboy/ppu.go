package gameboy

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type PPU struct {
	mapper *Mapper
	pallet [4]color.RGBA
	screen *ebiten.Image

	// Cache of all the 8x8 tiles, updated when the VRAM is written to
	tiles [128 * 3]*ebiten.Image

	// Cache of all the sprites, updated when the OAM is written to
	sprites [40]Sprite

	// Scanline register
	scanline   byte
	dotCounter int
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
		{0xFF, 0xFF, 0xFF, 0xFF}, // 0 is white
		{0x55, 0x55, 0x55, 0xFF}, // 1 is dark grey
		{0xAA, 0xAA, 0xAA, 0xFF}, // 2 is light grey
		{0x00, 0x00, 0x00, 0xFF}, // 3 is black
	}

	ppu := &PPU{
		pallet: pallet,
		mapper: mapper,

		// Internal screen buffer
		screen: ebiten.NewImage(160, 144),
	}

	// Horrible! Bidirectional dependency mess
	mapper.ppu = ppu

	// Empty tile cache
	for i := 0; i < 384; i++ {
		ppu.tiles[i] = ebiten.NewImage(8, 8)
	}

	// Empty sprite cache
	for i := 0; i < 40; i++ {
		ppu.sprites[i] = Sprite{}
	}

	return ppu
}

// Tile cache maps the tile data 8000-97FF to a cache of 384 8x8 images
func (ppu *PPU) updateTileCache(addr uint16) {
	tileNum := (addr - TILE_DATA_0) / 16
	tileDataAddr := TILE_DATA_0 + tileNum*16
	tileData := [16]byte{}
	for i := uint16(0); i < 16; i++ {
		tileData[i] = ppu.mapper.Read(tileDataAddr + i)
	}

	// This converts the 16 bytes of tile data into an 8x8 image
	// Using 2 bits per pixel to index the pallet
	img := ebiten.NewImage(8, 8)
	for tileByte := 0; tileByte < 16; tileByte += 2 {
		for bit := 0; bit < 8; bit++ {
			pIndex := ((tileData[tileByte] >> (7 - bit) & 1) << 1) | (tileData[tileByte+1] >> (7 - bit) & 1)
			colour := ppu.pallet[pIndex]
			img.Set(bit, tileByte/2, colour)
		}
	}

	// log.Printf("Updating tile cache for tile %d %04X", tileNum, tileDataAddr)

	ppu.tiles[tileNum] = img
}

func (ppu *PPU) updateSpriteCache(addr uint16) {
	i := (addr - OAM) / 4

	sprite := Sprite{
		y:        ppu.mapper.oam[i*4],
		x:        ppu.mapper.oam[i*4+1],
		tile:     ppu.mapper.oam[i*4+2],
		priority: ppu.mapper.oam[i*4+3]&0x80 == 0x80,
		flipY:    ppu.mapper.oam[i*4+3]&0x40 == 0x40,
		flipX:    ppu.mapper.oam[i*4+3]&0x20 == 0x20,
	}

	ppu.sprites[i] = sprite
}

func (ppu *PPU) render() {
	// TODO: Remove this later, it's helpful for debugging
	//randR := uint8(rand.Intn(255))
	//ppu.screen.Fill(color.RGBA{randR, 255, 255, 255})

	// Tile map is 9800-9BFF when LCDC bit 6 is NOT set
	tileMap := uint16(TILE_MAP_0)
	if ppu.GetLCDCBit(3) == 1 {
		// And 9C00-9FFF when it is set
		tileMap = uint16(TILE_MAP_1)
	}

	// Tile offset is 0x8800 when LCDC bit 4 is NOT set
	tileOffset := 256
	if ppu.GetLCDCBit(4) == 1 {
		// And 0x8000 when it is set
		tileOffset = 0
	}

	// get SCROLL_Y and SCROLL_X
	scrollY := ppu.mapper.Read(SCROLL_Y)
	scrollX := ppu.mapper.Read(SCROLL_X)

	// Read the 1024 bytes of tile map data
	// And render into the screen at the correct position
	for i := uint16(0); i < 1024; i++ {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64((i%32)*8), float64((i/32)*8))

		tilenum := int(ppu.mapper.Read(tileMap + i))
		if tileOffset > 0 {
			// Tile number is signed 8bit when LCDC bit 4 is NOT set
			tilenum = int(int8(tilenum))
		}
		op.GeoM.Translate(float64(-scrollX), float64(-scrollY))

		ppu.screen.DrawImage(ppu.tiles[tileOffset+tilenum], op)
	}

	// Handle OAM and render sprites
	for _, sprite := range ppu.sprites {
		if sprite.y == 0 && sprite.x == 0 || sprite.tile == 0 {
			continue
		}

		op := &ebiten.DrawImageOptions{}
		screenY := int(sprite.y) - 16
		screenX := int(sprite.x) - 8
		op.GeoM.Translate(float64(screenX), float64(screenY))
		ppu.screen.DrawImage(ppu.tiles[sprite.tile], op)
	}
}

func (ppu *PPU) cycle(clockCycles int) {
	ppu.dotCounter += clockCycles

	if ppu.dotCounter > 1 {
		ppu.dotCounter = 0

		ppu.scanline++

		if ppu.scanline > 153 {
			ppu.scanline = 0
			ppu.render()

			// request vblank interrupt
			ppu.mapper.requestInterrupt(0)
		}

		ppu.mapper.Write(0xFF44, ppu.scanline)
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
	return ppu.mapper.Read(LCD_CONTROL) >> bit & 1
}
