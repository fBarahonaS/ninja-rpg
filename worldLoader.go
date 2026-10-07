package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type ldtkTile struct {
	PX  [2]int `json:"px"`
	Src [2]int `json:"src"`
	T   int    `json:"t"`
}

type ldtkLayer struct {
	Width       int        `json:"__cWid"`
	Height      int        `json:"__cHei"`
	TilesetPath string     `json:"__tilesetRelPath"`
	IntGrid     []int      `json:"intGridCsv"`
	Tiles       []ldtkTile `json:"autoLayerTiles"`
}

type ldtkLevel struct {
	Identifier     string      `json:"identifier"`
	LayerInstances []ldtkLayer `json:"layerInstances"`
}

type ldtkRoot struct {
	Levels []ldtkLevel `json:"levels"`
}

// Own structures to hold the data we need for our game
type WorldTile struct {
	PosX, PosY float64
	SrcX, SrcY int
}

type World struct {
	Width       int
	Height      int
	IntGrid     []int
	Tiles       []WorldTile
	TilesetPath string
}

func LoadWorld(worldPath string) (*World, error) {
	worldData, err := os.ReadFile(worldPath)
	if err != nil {
		return nil, fmt.Errorf("load world %s: %w", worldPath, err)
	}

	var worldJson ldtkRoot
	err = json.Unmarshal(worldData, &worldJson)
	if err != nil {
		return nil, fmt.Errorf("load world %s: %w", worldPath, err)
	}

	if len(worldJson.Levels) == 0 {
		return nil, fmt.Errorf("no levels in %s", worldPath)
	}

	level := worldJson.Levels[0]

	if len(level.LayerInstances) == 0 {
		return nil, fmt.Errorf("no layers in %s", worldPath)
	}

	layer := level.LayerInstances[0]

	return &World{
		Width:       layer.Width,
		Height:      layer.Height,
		IntGrid:     layer.IntGrid,
		Tiles:       convertTiles(layer.Tiles),
		TilesetPath: filepath.Join(filepath.Dir(worldPath), layer.TilesetPath),
	}, nil
}

func convertTiles(layerTiles []ldtkTile) []WorldTile {
	// Reserving space up front
	// this avoids repeated allocations in append process
	tiles := make([]WorldTile, 0, len(layerTiles))

	for _, tile := range layerTiles {
		tiles = append(tiles, WorldTile{
			PosX: float64(tile.PX[0]),
			PosY: float64(tile.PX[1]),
			SrcX: tile.Src[0],
			SrcY: tile.Src[1],
		})
	}

	return tiles
}
