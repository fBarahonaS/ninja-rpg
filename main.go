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
	playerInitialSpeed     = 2.0

	// Just for now, we will use a fixed number of tiles per row.
	// In the future, we can calculate this based on the image size.
	tilesPerRow = 22
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
	player    *Player
	tileSet   *ebiten.Image
	levelSoil [][]int
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

func (g *Game) tileImgFromID(tileID int) *ebiten.Image {
	currentY := (tileID / tilesPerRow) * tileSize
	currentX := (tileID % tilesPerRow) * tileSize

	return g.tileSet.SubImage(
		image.Rect(currentX, currentY, currentX+tileSize, currentY+tileSize),
	).(*ebiten.Image)
}

func (g *Game) Draw(screen *ebiten.Image) {
	options := &ebiten.DrawImageOptions{}

	for row, cols := range g.levelSoil {
		for col, tileID := range cols {
			options.GeoM.Reset()
			options.GeoM.Translate(float64(col*tileSize), float64(row*tileSize))
			screen.DrawImage(
				g.tileImgFromID(tileID),
				options,
			)
		}
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

	// Loading images
	tileSetImg, _, err := ebitenutil.NewImageFromFile("assets/maps/tileSetFloor.png")
	if err != nil {
		handleFatalTermination(err)
	}

	playerImg, _, err := ebitenutil.NewImageFromFile("assets/img/ninja.png")
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
		levelSoil: [][]int{
			{466, 463, 463, 463, 463, 463, 463, 463, 463, 463, 463, 463, 463, 463, 463, 463, 463, 463, 463, 469},
			{484, 485, 485, 485, 485, 485, 485, 485, 485, 485, 485, 485, 485, 550, 550, 550, 550, 550, 550, 486},
			{484, 485, 485, 485, 485, 550, 550, 550, 550, 550, 550, 550, 550, 550, 550, 550, 550, 550, 550, 486},
			{484, 485, 485, 485, 485, 550, 550, 550, 550, 485, 485, 550, 550, 550, 550, 550, 550, 550, 550, 486},
			{484, 485, 485, 550, 550, 550, 550, 485, 485, 485, 485, 485, 485, 550, 550, 550, 550, 550, 550, 486},
			{484, 485, 485, 485, 550, 550, 550, 550, 550, 485, 485, 485, 485, 550, 550, 550, 550, 550, 550, 486},
			{484, 485, 485, 485, 550, 550, 550, 485, 485, 485, 485, 485, 485, 550, 550, 550, 550, 550, 550, 486},
			{484, 550, 550, 550, 550, 550, 550, 485, 485, 485, 485, 485, 485, 550, 550, 550, 550, 550, 550, 486},
			{484, 550, 550, 550, 550, 550, 550, 485, 485, 485, 485, 485, 485, 485, 485, 550, 550, 550, 485, 486},
			{484, 485, 550, 485, 550, 550, 550, 485, 485, 485, 485, 485, 485, 485, 485, 550, 550, 550, 485, 486},
			{484, 485, 550, 485, 550, 550, 550, 485, 485, 485, 485, 485, 485, 485, 485, 550, 550, 550, 485, 486},
			{484, 485, 550, 485, 550, 550, 550, 550, 550, 550, 485, 485, 485, 485, 485, 550, 550, 550, 485, 486},
			{488, 485, 550, 485, 485, 550, 550, 550, 550, 550, 550, 550, 485, 485, 485, 550, 485, 550, 485, 486},
			{510, 485, 485, 485, 485, 485, 485, 485, 485, 485, 485, 485, 485, 485, 485, 485, 485, 485, 485, 513},
			{532, 507, 507, 507, 507, 507, 507, 507, 507, 507, 507, 507, 507, 507, 507, 507, 507, 507, 507, 535},
		},
	}

	if err := ebiten.RunGame(gameInstance); err != nil {
		handleFatalTermination(err)
	}
}

func handleFatalTermination(err error) {
	log.Fatal(err)
}
