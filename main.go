package main

import (
	"image/color"
	"log"
	"os"
	"yarg/gameboy"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var (
	gb         *gameboy.Gameboy
	faceSource *text.GoTextFaceSource
)

const scale = 4

type Game struct{}

// Update the game state by one tick, happens at 2Mhz or two dots
func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		log.Println("Quitting...")
		os.Exit(0)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyL) {
		gameboy.SetLogging(true)
	}

	// Run the system for two dots (like a tick)
	gb.Cycle()

	return nil
}

func init() {
	fontFile, err := os.Open("fonts/hermit.otf")
	if err != nil {
		log.Fatal(err)
	}
	defer fontFile.Close()

	faceSource, err = text.NewGoTextFaceSource(fontFile)
	if err != nil {
		log.Fatal(err)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Main emulator screen
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.Filter = ebiten.FilterLinear
	screen.DrawImage(gb.GetScreen(), op)

	// Debug info
	msg := gb.GetDebugInfo()
	textOp := &text.DrawOptions{}
	textOp.GeoM.Translate(650, 20)
	textOp.LineSpacing = 22

	textOp.ColorScale.ScaleWithColor(color.RGBA{0x00, 0xee, 0x11, 0xff})
	text.Draw(screen, msg, &text.GoTextFace{
		Source: faceSource,
		Size:   20,
	}, textOp)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 160*4 + 500, 144 * 4
}

func main() {
	gb = gameboy.NewGameboy()
	gb.Start()

	//gb.LoadMemDump("roms/hello-world.dump")
	gb.LoadROM("roms/Tetris.gb")
	//gb.LoadROM("roms/hello-world.gb")

	game := &Game{}
	ebiten.SetTPS(4194304 / 2) // 4.194304MHz but tick the dots at half this speed
	ebiten.SetWindowSize(160*scale+500, 144*scale)
	ebiten.SetWindowTitle("Goboy Emulator (DMGO)")

	// Call ebiten.RunGame to start your game loop.
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
