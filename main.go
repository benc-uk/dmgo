package main

import (
	"dmgo/gameboy"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"

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
	scale      = 4
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

	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		gb.Buttons.Set("Right", true)
	} else if inpututil.IsKeyJustReleased(ebiten.KeyRight) {
		gb.Buttons.Set("Right", false)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		gb.Buttons.Set("Left", true)
	} else if inpututil.IsKeyJustReleased(ebiten.KeyLeft) {
		gb.Buttons.Set("Left", false)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		gb.Buttons.Set("Up", true)
	} else if inpututil.IsKeyJustReleased(ebiten.KeyUp) {
		gb.Buttons.Set("Up", false)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		gb.Buttons.Set("Down", true)
	} else if inpututil.IsKeyJustReleased(ebiten.KeyDown) {
		gb.Buttons.Set("Down", false)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		gb.Buttons.Set("A", true)
	} else if inpututil.IsKeyJustReleased(ebiten.KeyA) {
		gb.Buttons.Set("A", false)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		gb.Buttons.Set("B", true)
	} else if inpututil.IsKeyJustReleased(ebiten.KeyS) {
		gb.Buttons.Set("B", false)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		gb.Buttons.Set("Start", true)
	} else if inpututil.IsKeyJustReleased(ebiten.KeyEnter) {
		gb.Buttons.Set("Start", false)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		gb.Buttons.Set("Select", true)
	} else if inpututil.IsKeyJustReleased(ebiten.KeyBackspace) {
		gb.Buttons.Set("Select", false)
	}

	// Main emulator loop
	gb.Update(clockSpeed / 20)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Render emulator screen
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(scale), float64(scale))
	op.Filter = ebiten.FilterNearest
	screen.DrawImage(gb.GetScreen(), op)

	// Debug info
	msg := gb.GetDebugInfo()
	textOp := &text.DrawOptions{}
	textOp.GeoM.Translate(float64(163*scale), 20)
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
	if len(os.Args) > 1 {
		gb.LoadROM(os.Args[1])
	} else {
		log.Println("No game cart ROM specified, booting without a cart")
	}

	gb.Running = true

	game := &Game{}
	ebiten.SetWindowSize(160*scale+140*scale, 144*scale)
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
