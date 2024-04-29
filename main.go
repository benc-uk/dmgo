package main

import (
	"image"
	"image/color"
	"image/png"
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

const (
	scale      = 4
	clockSpeed = 4194304
)

type Game struct{}

func init() {
	// Load font
	fontFile, err := os.Open("res/hermit.otf")
	if err != nil {
		log.Fatal(err)
	}
	defer fontFile.Close()

	faceSource, err = text.NewGoTextFaceSource(fontFile)
	if err != nil {
		log.Fatal(err)
	}

	// Load icon
	iconFile, err := os.Open("res/icon.png")
	if err != nil {
		log.Fatal(err)
	}
	defer iconFile.Close()

	icon, err := png.Decode(iconFile)
	if err != nil {
		log.Fatal(err)
	}

	ebiten.SetWindowIcon([]image.Image{icon})
}

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

// Entry point is here
func main() {
	gb = gameboy.NewGameboy()
	gb.Start()

	//gb.LoadMemDump("roms/hello-world.dump")
	//gb.LoadROM("roms/cpu_instrs.gb")
	//gb.LoadROM("roms/hello-world.gb")
	gb.LoadROM("roms/Tetris.gb")

	game := &Game{}
	ebiten.SetTPS(clockSpeed / 2) // 4.19MHz but we run at half this for the PPU and CPU
	ebiten.SetWindowSize(160*scale+500, 144*scale)
	ebiten.SetWindowTitle("Gameboy Emulator (DMGO)")

	// Call ebiten.RunGame to start
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
