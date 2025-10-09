package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type TileSet struct {
	Columns      int    `json:"columns"`
	Image        string `json:"image"`
	ImageHeight  int    `json:"imageheight"`
	ImageWidth   int    `json:"imagewidth"`
	Margin       int    `json:"margin"`
	Name         string `json:"name"`
	Spacing      int    `json:"spacing"`
	TileCount    int    `json:"tilecount"`
	TileHeight   int    `json:"tileheight"`
	TileWidth    int    `json:"tilewidth"`
	Type         string `json:"type"`
	Version      string `json:"version"`
	TiledVersion string `json:"tiledversion"`
	Texture      rl.Texture2D
}

type TileMap struct {
	CompressionLevel int            `json:"compressionlevel"`
	Width            int            `json:"width"`
	Height           int            `json:"height"`
	Infinite         bool           `json:"infinite"`
	Layers           []TileMapLayer `json:"layers"`
	NextLayerID      int            `json:"nextlayerid"`
	NextObjectID     int            `json:"nextobjectid"`
	Orientation      string         `json:"orientation"`
	RenderOrder      string         `json:"renderorder"`
	TileHeight       int            `json:"tileheight"`
	TileWidth        int            `json:"tilewidth"`
	Type             string         `json:"type"`
	TiledVersion     string         `json:"tiledversion"`
	Version          string         `json:"version"`
	TileSets         []struct {
		FirstGID int    `json:"firstgid"`
		Source   string `json:"source"`
	} `json:"tilesets"`
}

type TileMapLayer struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Data    []int   `json:"data"`
	Width   int     `json:"width"`
	Height  int     `json:"height"`
	Opacity float32 `json:"opacity"`
	Type    string  `json:"type"`
	Visible bool    `json:"visible"`
	X       int     `json:"x"`
	Y       int     `json:"y"`
}

func LoadTileSet(filename string) *TileSet {

	tileset := &TileSet{}

	jsonFile, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, tileset)
	if err != nil {
		panic(err)
	}

	return tileset
}

func LoadTileMap(filename string) *TileMap {

	tilemap := &TileMap{}

	jsonFile, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, tilemap)
	if err != nil {
		panic(err)
	}

	return tilemap
}
