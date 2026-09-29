package main

import (
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	gameTitle    = ":: Ninja Hero RPG ::"
	screenWidth  = 640
	screenHeight = 480
	tileSize     = 16
)

var (
	// Defining variable to handle player position
	PlayerPosX  float64 = 100
	PlayerPosY  float64 = 100
	playerSpeed float64 = 4.0
)

type Game struct {
	// Defining variable to handle player image
	PlayerImg *ebiten.Image
}

func (g *Game) Update() error {
	// Reacting to pressed keys
	switch {
	case ebiten.IsKeyPressed(ebiten.KeyArrowUp):
		if PlayerPosY > 0 {
			PlayerPosY -= playerSpeed
		}
	case ebiten.IsKeyPressed(ebiten.KeyArrowDown):
		if PlayerPosY < screenHeight-tileSize {
			PlayerPosY += playerSpeed
		}
	case ebiten.IsKeyPressed(ebiten.KeyArrowLeft):
		if PlayerPosX > 0 {
			PlayerPosX -= playerSpeed
		}
	case ebiten.IsKeyPressed(ebiten.KeyArrowRight):
		if PlayerPosX < screenWidth-tileSize {
			PlayerPosX += playerSpeed
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{121, 180, 255, 255})

	options := &ebiten.DrawImageOptions{}
	options.GeoM.Translate(PlayerPosX, PlayerPosY)

	// Drawing the player image within window
	screen.DrawImage(
		g.PlayerImg.SubImage(
			image.Rect(0, 0, tileSize, tileSize),
		).(*ebiten.Image),
		options,
	)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ebiten.WindowSize()
}

func main() {
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle(gameTitle)

	// Loading the player image
	PlayerImg, _, err := ebitenutil.NewImageFromFile("assets/img/ninja.png")
	if err != nil {
		handleFatalTermination(err)
	}

	// Implementing the game instance
	gameInstance := &Game{
		PlayerImg: PlayerImg,
	}

	if err := ebiten.RunGame(gameInstance); err != nil {
		handleFatalTermination(err)
	}
}

func handleFatalTermination(err error) {
	log.Fatal(err)
}
