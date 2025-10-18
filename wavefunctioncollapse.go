package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// Tile Edges
var tileEdges = map[string]int{
	"GRASS":    0,
	"WATER":    1,
	"FOREST":   2,
	"COAST_N":  3,
	"COAST_E":  4,
	"COAST_S":  5,
	"COAST_W":  6,
	"FOREST_N": 7,
	"FOREST_E": 8,
	"FOREST_S": 9,
	"FOREST_W": 10,
	"ROCK_N":   11,
	"ROCK_E":   12,
	"ROCK_S":   13,
	"ROCK_W":   14,
	"ROCK":     15,
}

type WaveParams struct {
	Data []WaveTile `json:"data"`
}

type WaveTile struct {
	ID     int      `json:"id"`
	Type   string   `json:"type"`
	Desc   []string `json:"desc"`
	Weight int      `json:"weight"`
}

var waveParams = &WaveParams{}

func WaveLoadParams(filename string) *WaveParams {

	params := &WaveParams{}
	jsonFile, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, params)
	if err != nil {
		panic(err)
	}

	return params
}

func WaveInit(tileSet *TileSet) {

	waveParams = WaveLoadParams("waveparams.json")
	WaveGetTileIds(tileSet, waveParams)
}

func WaveGetTileIds(tileset *TileSet, waveParams *WaveParams) {

	mapIds := map[string]int{}
	for _, tile := range tileset.Tiles {
		if tile.Type != "" {
			mapIds[tile.Type] = tile.ID
		}
	}

	for _, tile := range waveParams.Data {
		tile.ID = mapIds[tile.Type]
	}
}

func WaveGenerate(columns int, rows int) *[]int {

	// grid := make([]int, columns*rows)

	// r := rand.New(rand.NewSource(time.Now().UnixNano()))
	// i := r.Intn(columns * rows)

	// cell_n := -1
	// if i > columns {
	// 	cell_n = grid[i-columns]
	// }
	// cell_e := -1
	// if (i+1)%columns != 0 {
	// 	cell_e = grid[i+1]
	// }
	// cell_s := -1
	// if i < (rows-1)*columns {
	// 	cell_s = grid[i+columns]
	// }
	// cell_w := -1
	// if i%columns != 0 {
	// 	cell_w = grid[i-1]
	// }

	// possible_tiles := make([]int, 0, 32)
	// possible_weights := make([]int, 0, 32)
	// total_weight := 0

	// for tile_id, rules := range tileRules {
	// 	fmt.Println(tile_id, rules)
	// 	break
	// }

	// return &grid
	return nil
}
