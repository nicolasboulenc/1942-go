package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Directions
const (
	NORTH = 0
	EAST  = 1
	SOUTH = 2
	WEST  = 3
)

// Tile Edges
const (
	GRASS    = 0
	WATER    = 1
	FOREST   = 2
	COAST_N  = 3
	COAST_E  = 4
	COAST_S  = 5
	COAST_W  = 6
	FOREST_N = 7
	FOREST_E = 8
	FOREST_S = 9
	FOREST_W = 10
	ROCK_N   = 11
	ROCK_E   = 12
	ROCK_S   = 13
	ROCK_W   = 14
	ROCK     = 15
)

// Tile Types
type TileIds struct {
	TILE_GRASS      int
	TILE_WATER      int
	TILE_FOREST     int
	TILE_COAST_N    int
	TILE_COAST_E    int
	TILE_COAST_S    int
	TILE_COAST_W    int
	TILE_COAST_NE   int
	TILE_COAST_SE   int
	TILE_COAST_SW   int
	TILE_COAST_NW   int
	TILE_COAST_NE2  int
	TILE_COAST_SE2  int
	TILE_COAST_SW2  int
	TILE_COAST_NW2  int
	TILE_ROCK_N     int
	TILE_ROCK_E     int
	TILE_ROCK_S     int
	TILE_ROCK_W     int
	TILE_ROCK_NE    int
	TILE_ROCK_SE    int
	TILE_ROCK_SW    int
	TILE_ROCK_NW    int
	TILE_FOREST_N   int
	TILE_FOREST_E   int
	TILE_FOREST_S   int
	TILE_FOREST_W   int
	TILE_FOREST_NE  int
	TILE_FOREST_SE  int
	TILE_FOREST_SW  int
	TILE_FOREST_NW  int
	TILE_FOREST_NE2 int
	TILE_FOREST_SE2 int
	TILE_FOREST_SW2 int
	TILE_FOREST_NW2 int
}

// Dictionary of all tile types and tile edges, on the directions [North, East, South, West]
type TileRules struct {
	TILE_GRASS      [4]int
	TILE_WATER      [4]int
	TILE_FOREST     [4]int
	TILE_COAST_N    [4]int
	TILE_COAST_E    [4]int
	TILE_COAST_S    [4]int
	TILE_COAST_W    [4]int
	TILE_COAST_NE   [4]int
	TILE_COAST_SE   [4]int
	TILE_COAST_SW   [4]int
	TILE_COAST_NW   [4]int
	TILE_COAST_NE2  [4]int
	TILE_COAST_SE2  [4]int
	TILE_COAST_SW2  [4]int
	TILE_COAST_NW2  [4]int
	TILE_ROCK_N     [4]int
	TILE_ROCK_E     [4]int
	TILE_ROCK_S     [4]int
	TILE_ROCK_W     [4]int
	TILE_ROCK_NE    [4]int
	TILE_ROCK_SE    [4]int
	TILE_ROCK_SW    [4]int
	TILE_ROCK_NW    [4]int
	TILE_FOREST_N   [4]int
	TILE_FOREST_E   [4]int
	TILE_FOREST_S   [4]int
	TILE_FOREST_W   [4]int
	TILE_FOREST_NE  [4]int
	TILE_FOREST_SE  [4]int
	TILE_FOREST_SW  [4]int
	TILE_FOREST_NW  [4]int
	TILE_FOREST_NE2 [4]int
	TILE_FOREST_SE2 [4]int
	TILE_FOREST_SW2 [4]int
	TILE_FOREST_NW2 [4]int
}

type TileWeights struct {
	TILE_GRASS      int
	TILE_WATER      int
	TILE_FOREST     int
	TILE_COAST_N    int
	TILE_COAST_E    int
	TILE_COAST_S    int
	TILE_COAST_W    int
	TILE_COAST_NE   int
	TILE_COAST_SE   int
	TILE_COAST_SW   int
	TILE_COAST_NW   int
	TILE_COAST_NE2  int
	TILE_COAST_SE2  int
	TILE_COAST_SW2  int
	TILE_COAST_NW2  int
	TILE_FOREST_N   int
	TILE_FOREST_E   int
	TILE_FOREST_S   int
	TILE_FOREST_W   int
	TILE_FOREST_NE  int
	TILE_FOREST_SE  int
	TILE_FOREST_SW  int
	TILE_FOREST_NW  int
	TILE_FOREST_NE2 int
	TILE_FOREST_SE2 int
	TILE_FOREST_SW2 int
	TILE_FOREST_NW2 int
	TILE_ROCK_N     int
	TILE_ROCK_E     int
	TILE_ROCK_S     int
	TILE_ROCK_W     int
	TILE_ROCK_NE    int
	TILE_ROCK_SE    int
	TILE_ROCK_SW    int
	TILE_ROCK_NW    int
}

var tileIds = TileIds{}

var tileRules = TileRules{
	TILE_GRASS:      [4]int{GRASS, GRASS, GRASS, GRASS},
	TILE_WATER:      [4]int{WATER, WATER, WATER, WATER},
	TILE_FOREST:     [4]int{FOREST, FOREST, FOREST, FOREST},
	TILE_COAST_N:    [4]int{GRASS, COAST_N, WATER, COAST_N},
	TILE_COAST_E:    [4]int{COAST_E, GRASS, COAST_E, WATER},
	TILE_COAST_S:    [4]int{WATER, COAST_S, GRASS, COAST_S},
	TILE_COAST_W:    [4]int{COAST_W, WATER, COAST_W, GRASS},
	TILE_COAST_NE:   [4]int{GRASS, GRASS, COAST_E, COAST_N},
	TILE_COAST_SE:   [4]int{COAST_E, GRASS, GRASS, COAST_S},
	TILE_COAST_SW:   [4]int{COAST_W, COAST_S, GRASS, GRASS},
	TILE_COAST_NW:   [4]int{GRASS, COAST_N, COAST_W, GRASS},
	TILE_COAST_NE2:  [4]int{COAST_E, COAST_N, WATER, WATER},
	TILE_COAST_SE2:  [4]int{WATER, COAST_S, COAST_E, WATER},
	TILE_COAST_SW2:  [4]int{WATER, WATER, COAST_W, COAST_S},
	TILE_COAST_NW2:  [4]int{COAST_W, WATER, WATER, COAST_N},
	TILE_ROCK_N:     [4]int{ROCK, ROCK_N, GRASS, ROCK_N},
	TILE_ROCK_E:     [4]int{ROCK_E, ROCK, ROCK_E, GRASS},
	TILE_ROCK_S:     [4]int{GRASS, ROCK_S, ROCK, ROCK_S},
	TILE_ROCK_W:     [4]int{ROCK_W, GRASS, ROCK_W, ROCK},
	TILE_ROCK_NE:    [4]int{ROCK_E, ROCK_N, GRASS, GRASS},
	TILE_ROCK_SE:    [4]int{GRASS, ROCK_S, ROCK_E, GRASS},
	TILE_ROCK_SW:    [4]int{GRASS, GRASS, ROCK_W, ROCK_S},
	TILE_ROCK_NW:    [4]int{ROCK_W, GRASS, GRASS, ROCK_N},
	TILE_FOREST_N:   [4]int{FOREST, FOREST_N, GRASS, FOREST_N},
	TILE_FOREST_E:   [4]int{FOREST_E, FOREST, FOREST_E, GRASS},
	TILE_FOREST_S:   [4]int{GRASS, FOREST_S, FOREST, FOREST_S},
	TILE_FOREST_W:   [4]int{FOREST_W, GRASS, FOREST_W, FOREST},
	TILE_FOREST_NE:  [4]int{FOREST_E, FOREST_N, GRASS, GRASS},
	TILE_FOREST_SE:  [4]int{GRASS, FOREST_S, FOREST_E, GRASS},
	TILE_FOREST_SW:  [4]int{GRASS, GRASS, FOREST_W, FOREST_S},
	TILE_FOREST_NW:  [4]int{FOREST_W, GRASS, GRASS, FOREST_N},
	TILE_FOREST_NE2: [4]int{FOREST, FOREST, FOREST_E, FOREST_N},
	TILE_FOREST_SE2: [4]int{FOREST_E, FOREST, FOREST, FOREST_S},
	TILE_FOREST_SW2: [4]int{FOREST_W, FOREST_S, FOREST, FOREST},
	TILE_FOREST_NW2: [4]int{FOREST, FOREST_N, FOREST_W, FOREST},
}

var tileWeights = TileWeights{
	TILE_GRASS:      16,
	TILE_WATER:      4,
	TILE_FOREST:     5,
	TILE_COAST_N:    5,
	TILE_COAST_E:    5,
	TILE_COAST_S:    5,
	TILE_COAST_W:    5,
	TILE_COAST_NE:   5,
	TILE_COAST_SE:   5,
	TILE_COAST_SW:   5,
	TILE_COAST_NW:   5,
	TILE_COAST_NE2:  2,
	TILE_COAST_SE2:  2,
	TILE_COAST_SW2:  2,
	TILE_COAST_NW2:  2,
	TILE_FOREST_N:   4,
	TILE_FOREST_E:   4,
	TILE_FOREST_S:   4,
	TILE_FOREST_W:   4,
	TILE_FOREST_NE:  4,
	TILE_FOREST_SE:  4,
	TILE_FOREST_SW:  4,
	TILE_FOREST_NW:  4,
	TILE_FOREST_NE2: 2,
	TILE_FOREST_SE2: 2,
	TILE_FOREST_SW2: 2,
	TILE_FOREST_NW2: 2,
	TILE_ROCK_N:     4,
	TILE_ROCK_E:     4,
	TILE_ROCK_S:     4,
	TILE_ROCK_W:     4,
	TILE_ROCK_NE:    4,
	TILE_ROCK_SE:    4,
	TILE_ROCK_SW:    4,
	TILE_ROCK_NW:    4,
}

func WaveGetTileIds(tileset *TileSet, tileids *TileIds) {
	for i, tile := range tileset.Tiles {
		switch tile.Type {
		case "TILE_GRASS":
			tileids.TILE_GRASS = tileset.Tiles[i].ID
		case "TILE_WATER":
			tileids.TILE_WATER = tileset.Tiles[i].ID
		case "TILE_FOREST":
			tileids.TILE_FOREST = tileset.Tiles[i].ID
		case "TILE_COAST_N":
			tileids.TILE_COAST_N = tileset.Tiles[i].ID
		case "TILE_COAST_E":
			tileids.TILE_COAST_E = tileset.Tiles[i].ID
		case "TILE_COAST_S":
			tileids.TILE_COAST_S = tileset.Tiles[i].ID
		case "TILE_COAST_W":
			tileids.TILE_COAST_W = tileset.Tiles[i].ID
		case "TILE_COAST_NE":
			tileids.TILE_COAST_NE = tileset.Tiles[i].ID
		case "TILE_COAST_SE":
			tileids.TILE_COAST_SE = tileset.Tiles[i].ID
		case "TILE_COAST_SW":
			tileids.TILE_COAST_SW = tileset.Tiles[i].ID
		case "TILE_COAST_NW":
			tileids.TILE_COAST_NW = tileset.Tiles[i].ID
		case "TILE_COAST_NE2":
			tileids.TILE_COAST_NE2 = tileset.Tiles[i].ID
		case "TILE_COAST_SE2":
			tileids.TILE_COAST_SE2 = tileset.Tiles[i].ID
		case "TILE_COAST_SW2":
			tileids.TILE_COAST_SW2 = tileset.Tiles[i].ID
		case "TILE_COAST_NW2":
			tileids.TILE_COAST_NW2 = tileset.Tiles[i].ID
		case "TILE_ROCK_N":
			tileids.TILE_ROCK_N = tileset.Tiles[i].ID
		case "TILE_ROCK_E":
			tileids.TILE_ROCK_E = tileset.Tiles[i].ID
		case "TILE_ROCK_S":
			tileids.TILE_ROCK_S = tileset.Tiles[i].ID
		case "TILE_ROCK_W":
			tileids.TILE_ROCK_W = tileset.Tiles[i].ID
		case "TILE_ROCK_NE":
			tileids.TILE_ROCK_NE = tileset.Tiles[i].ID
		case "TILE_ROCK_SE":
			tileids.TILE_ROCK_SE = tileset.Tiles[i].ID
		case "TILE_ROCK_SW":
			tileids.TILE_ROCK_SW = tileset.Tiles[i].ID
		case "TILE_ROCK_NW":
			tileids.TILE_ROCK_NW = tileset.Tiles[i].ID
		case "TILE_FOREST_N":
			tileids.TILE_FOREST_N = tileset.Tiles[i].ID
		case "TILE_FOREST_E":
			tileids.TILE_FOREST_E = tileset.Tiles[i].ID
		case "TILE_FOREST_S":
			tileids.TILE_FOREST_S = tileset.Tiles[i].ID
		case "TILE_FOREST_W":
			tileids.TILE_FOREST_W = tileset.Tiles[i].ID
		case "TILE_FOREST_NE":
			tileids.TILE_FOREST_NE = tileset.Tiles[i].ID
		case "TILE_FOREST_SE":
			tileids.TILE_FOREST_SE = tileset.Tiles[i].ID
		case "TILE_FOREST_SW":
			tileids.TILE_FOREST_SW = tileset.Tiles[i].ID
		case "TILE_FOREST_NW":
			tileids.TILE_FOREST_NW = tileset.Tiles[i].ID
		case "TILE_FOREST_NE2":
			tileids.TILE_FOREST_NE2 = tileset.Tiles[i].ID
		case "TILE_FOREST_SE2":
			tileids.TILE_FOREST_SE2 = tileset.Tiles[i].ID
		case "TILE_FOREST_SW2":
			tileids.TILE_FOREST_SW2 = tileset.Tiles[i].ID
		case "TILE_FOREST_NW2":
			tileids.TILE_FOREST_NW2 = tileset.Tiles[i].ID
		}
	}
}

func WaveInit(tileSet *TileSet) {

	WaveGetTileIds(tileSet, &tileIds)
}

func WaveGenerate(columns int, rows int) *[]int {

	grid := make([]int, columns*rows)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	i := r.Intn(columns * rows)

	cell_n := -1
	if i > columns {
		cell_n = grid[i-columns]
	}
	cell_e := -1
	if (i+1)%columns != 0 {
		cell_e = grid[i+1]
	}
	cell_s := -1
	if i < (rows-1)*columns {
		cell_s = grid[i+columns]
	}
	cell_w := -1
	if i%columns != 0 {
		cell_w = grid[i-1]
	}

	possible_tiles := make([]int, 0, 32)
	possible_weights := make([]int, 0, 32)
	total_weight := 0

	for tile_id, rules := range tileRules {
		fmt.Println(tile_id, rules)
		break
	}

	return &grid
}
