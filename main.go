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

type Game struct {
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		gb.Tick()
	}

	return nil
}

func init() {
	f, err := os.Open("fonts/Crisp.ttf")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	s, err := text.NewGoTextFaceSource(f)
	if err != nil {
		log.Fatal(err)
	}
	faceSource = s
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.Filter = ebiten.FilterLinear

	screen.DrawImage(gb.GetScreen(), op)

	msg := gb.GetDebugInfo()
	textOp := &text.DrawOptions{}
	textOp.GeoM.Translate(650, 20)
	textOp.LineSpacing = 28.0

	textOp.ColorScale.ScaleWithColor(color.RGBA{0x00, 0xff, 0x11, 0xff})
	text.Draw(screen, msg, &text.GoTextFace{
		Source: faceSource,
		Size:   26,
	}, textOp)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 160*4 + 500, 144 * 4
}

func main() {
	gb = gameboy.NewGameboy()
	//gb.LoadMemDump("roms/Tetris.dmp")
	gb.Start()

	//gb.LoadROM("roms/Tetris.gb")
	gb.LoadROM("roms/hello-world.gb")
	// for i := 0; i < 1000000; i++ {
	// 	gb.Tick()
	// }

	game := &Game{}
	ebiten.SetWindowSize(160*scale+500, 144*scale)
	ebiten.SetWindowTitle("Gameboy Emulator (YARG)")

	// Call ebiten.RunGame to start your game loop.
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
