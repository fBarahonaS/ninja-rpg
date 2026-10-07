package main

import (
	"image"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	gameTitle              = ":: Ninja Hero RPG ::"
	windowWidth            = 640
	windowHeight           = 480
	canvasWidth            = 320
	canvasHeight           = 240
	tileSize               = 16
	canvasMiddleX          = canvasWidth / 2.0
	canvasMiddleY          = canvasHeight / 2.0
	tileMiddle             = tileSize / 2.0
	playerInitialPositionX = canvasMiddleX - tileMiddle
	playerInitialPositionY = canvasMiddleY - tileMiddle
	worldFilePath          = "assets/maps/world1.ldtk"
	playerSpriteImgPath    = "assets/img/ninja.png"
	playerInitialSpeed     = 2.0
)

type Sprite struct {
	img       *ebiten.Image
	positionX float64
	positionY float64
}

type Player struct {
	*Sprite
	speed float64
}

type Game struct {
	player  *Player
	tileSet *ebiten.Image
	world   *World
}

func (g *Game) Update() error {
	switch {
	case ebiten.IsKeyPressed(ebiten.KeyArrowUp):
		g.player.positionY = max(
			g.player.positionY-g.player.speed,
			0,
		)
	case ebiten.IsKeyPressed(ebiten.KeyArrowDown):
		g.player.positionY = min(
			g.player.positionY+g.player.speed,
			canvasHeight-tileSize,
		)
	case ebiten.IsKeyPressed(ebiten.KeyArrowLeft):
		g.player.positionX = max(
			g.player.positionX-g.player.speed,
			0,
		)
	case ebiten.IsKeyPressed(ebiten.KeyArrowRight):
		g.player.positionX = min(
			g.player.positionX+g.player.speed,
			canvasWidth-tileSize,
		)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	options := &ebiten.DrawImageOptions{}

	for _, tile := range g.world.Tiles {
		options.GeoM.Reset()
		options.GeoM.Translate(tile.PosX, tile.PosY)
		screen.DrawImage(
			g.tileSet.SubImage(
				image.Rect(tile.SrcX, tile.SrcY, tile.SrcX+tileSize, tile.SrcY+tileSize),
			).(*ebiten.Image),
			options,
		)
	}

	options = &ebiten.DrawImageOptions{}
	options.GeoM.Translate(g.player.positionX, g.player.positionY)
	screen.DrawImage(
		g.player.img.SubImage(
			image.Rect(0, 0, tileSize, tileSize),
		).(*ebiten.Image),
		options,
	)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return canvasWidth, canvasHeight
}

func main() {
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle(gameTitle)

	world, err := LoadWorld(worldFilePath)
	if err != nil {
		handleFatalTermination(err)
	}

	// Loading images
	tileSetImg, _, err := ebitenutil.NewImageFromFile(world.TilesetPath)
	if err != nil {
		handleFatalTermination(err)
	}

	playerImg, _, err := ebitenutil.NewImageFromFile(playerSpriteImgPath)
	if err != nil {
		handleFatalTermination(err)
	}

	// Game instance
	gameInstance := &Game{
		player: &Player{
			Sprite: &Sprite{
				img:       playerImg,
				positionX: playerInitialPositionX,
				positionY: playerInitialPositionY,
			},
			speed: playerInitialSpeed,
		},
		tileSet: tileSetImg,
		world:   world,
	}

	if err := ebiten.RunGame(gameInstance); err != nil {
		handleFatalTermination(err)
	}
}

func handleFatalTermination(err error) {
	log.Fatal(err)
}
