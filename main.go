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
	"gopkg.in/yaml.v2"
)

var (
	gb         *gameboy.Gameboy
	faceSource *text.GoTextFaceSource
	config     gameboy.Config
)

const (
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

// Update the game state by the given delta time
func (g *Game) Update() error {
	// Check for step mode
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) && !gb.Running {
		// Special case for telling the CPU to step
		gb.Update(-1)
	}

	tps := int(ebiten.ActualTPS())
	if tps <= 0 {
		return nil
	}

	// Main emulator loop
	gb.Update(clockSpeed / tps)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Render emulator screen
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(config.Scale), float64(config.Scale))
	op.Filter = ebiten.FilterLinear
	screen.DrawImage(gb.GetScreen(), op)

	// Debug info
	msg := gb.GetDebugInfo()
	textOp := &text.DrawOptions{}
	textOp.GeoM.Translate(float64(163*config.Scale), 20)
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
	// Read config.yaml file
	configFile, err := os.Open("./config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	defer configFile.Close()

	config, err = readConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}

	gb = gameboy.NewGameboy(config)
	if config.ROM != "" {
		gb.LoadROM(config.ROM)
	} else {
		log.Println("No game cart ROM specified, booting without a cart")
	}
	gb.Running = true

	game := &Game{}
	ebiten.SetWindowSize(160*config.Scale+140*config.Scale, 144*config.Scale)
	ebiten.SetWindowTitle("Gameboy Emulator (DMGO)")

	// Call ebiten.RunGame to start
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func readConfig(file *os.File) (gameboy.Config, error) {
	// Read the file
	decoder := yaml.NewDecoder(file)
	err := decoder.Decode(&config)
	if err != nil {
		return gameboy.Config{}, err
	}

	return config, nil
}
