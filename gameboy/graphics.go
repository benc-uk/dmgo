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
	tiles [256]*ebiten.Image

	// Cache of all the sprites, updated when the OAM is written to
	sprites [40]Sprite

	// Scanline register
	scanline byte
}

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
	for i := 0; i < 256; i++ {
		ppu.tiles[i] = ebiten.NewImage(8, 8)
	}

	// Empty sprite cache
	for i := 0; i < 40; i++ {
		ppu.sprites[i] = Sprite{}
	}

	return ppu
}

func (ppu *PPU) updateTileCache(addr uint16) {
	tileNum := (addr - TILE_DATA_1) / 16
	tileDataAddr := TILE_DATA_1 + tileNum*16
	tileData := [16]byte{}
	for i := uint16(0); i < 16; i++ {
		tileData[i] = ppu.mapper.Read(tileDataAddr + i)
	}

	img := ebiten.NewImage(8, 8)
	for tileByte := 0; tileByte < 16; tileByte += 2 {
		for bit := 0; bit < 8; bit++ {
			pIndex := ((tileData[tileByte] >> (7 - bit) & 1) << 1) | (tileData[tileByte+1] >> (7 - bit) & 1)
			colour := ppu.pallet[pIndex]
			img.Set(bit, tileByte/2, colour)
		}
	}

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

func (ppu *PPU) Render() {
	// TODO: REMOVE THIS?
	ppu.screen.Fill(color.RGBA{191, 232, 183, 255})

	tileMap := uint16(TILE_MAP_1)

	// Read the 1024 bytes of tile map data
	// And render into the screen at the correct position
	for i := uint16(0); i < 1024; i++ {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64((i%32)*8), float64((i/32)*8))
		ppu.screen.DrawImage(ppu.tiles[ppu.mapper.Read(tileMap+i)], op)
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

	// HACK: This is a hack to update the scanline register
	ppu.scanline++
	if ppu.scanline > 153 {
		ppu.scanline = 0
	}

	// Write the scanline register to memory
	ppu.mapper.Write(0xFF44, ppu.scanline)
}
