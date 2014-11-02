package main

import (
	"fmt"
	term "github.com/buger/goterm"
	"math/rand"
	"time"
	"flag"
)

type Tile int64

const (
	WORLD_X = 200
	WORLD_Y = 200
	WORLD_SCALE = 2
)
const (
	TILE_EMPTY = 1 << iota
	TILE_GRASS = 1 << iota
	TILE_DIRT  = 1 << iota
	TILE_LAVA  = 1 << iota
	TILE_WATER = 1 << iota
	TILE_LAND  = TILE_GRASS | TILE_DIRT
	TILE_FLUID = TILE_LAVA | TILE_WATER
)

type TileRow []Tile
type Layer []TileRow
type Layers []Layer
type World struct {
	Layers Layers
	Width int
	Height int
}

func (t Tile) String() string {
	switch t {
	case TILE_EMPTY:
		return term.Color("#", term.BLACK)
	case TILE_GRASS:
		return term.Color("G", term.GREEN)
	case TILE_DIRT:
		return term.Color("D", term.YELLOW)
	case TILE_LAVA:
		return term.Color("L", term.RED)
	case TILE_WATER:
		return term.Color("W", term.BLUE)
	default:
		return "?"
	}
}

func (tr TileRow) String() string {
	row := ""
	for _, tile := range tr {
		for i := 0; i < WORLD_SCALE; i++ {
			row += tile.String() 
		}
	}
	return row
}

func (l Layer) String() string {
	layer := ""
	for _, row := range l {
		for i := 0; i < WORLD_SCALE; i++ {
			layer += row.String() + "\n"
		}
	}
	return layer
}

func (l Layer) IsTile(x, y int, tile Tile) bool {
	
	if ((x >= 0 && y >= 0) && y < len(l)) {
		if x < len(l[y]) {
			return l[y][x] & tile != 0
		}
	}
	return false;
}

func (l Layer) ClosestLand(x, y int) Tile {
	if l.IsTile(x,y,TILE_LAND) {
		return l[y][x]
	} else {
		if l.IsTile(x-1,y, TILE_LAND) {
			return l[y][x-1]
		} else if l.IsTile(x,y-1, TILE_LAND) {
			return l[y-1][x]
		} else if l.IsTile(x-1,y-1, TILE_LAND) {
			return l[y-1][x-1]
		} else if l.IsTile(x+1,y, TILE_LAND) {
			return l[y][x+1]
		} else if l.IsTile(x,y+1, TILE_LAND) {
			return l[y+1][x]
		} else if l.IsTile(x+1,y+1, TILE_LAND) {
			return l[y+1][x+1]
		} else if l.IsTile(x-1,y+1, TILE_LAND) {
			return l[y+1][x-1]
		} else if l.IsTile(x+1,y-1, TILE_LAND) {
			return l[y-1][x+1]
		} else {
			return TILE_EMPTY
		}
	}
}

func (w World) String() string {
	world := ""
	for z, layer := range w.Layers {
		world += fmt.Sprintf("Layer: %d\n", z)
		world += layer.String() + "\n"
	}
	return world
}

func (w *World) Scale(factor int) *World {
	height := w.Height
	width := w.Width

	world := World{}
	world.Layers = make(Layers, len(w.Layers))
	world.Width = width * factor
	world.Height = height * factor
	layers := world.Layers
	for z := 0; z < len(w.Layers); z++ {
		layers[z] = make(Layer, world.Height)
		for y := 0; y < height; y++ {
			for yFactor := 0; yFactor < factor; yFactor++ {
				layers[z][(y*factor)+yFactor] = make(TileRow, world.Width)
				for x := 0; x < width; x++ {
					for xFactor := 0; xFactor < factor; xFactor++ {
						layers[z][(y*factor)+yFactor][(x*factor)+xFactor] = w.Layers[z][y][x]
					}
				}
			}
		}
	}

	return &world
}

type WorldSeed map[Tile]int

func (ws WorldSeed) NextTile() Tile {
	rand.Seed(time.Now().UTC().UnixNano())

	sum := 0
	for _, weight := range ws {
		sum += weight
	}

	accum := 0
	num := rand.Intn(sum)
	for tile, weight := range ws {
		if weight != 0 &&
			num >= accum &&
			num <= (accum+weight) {
			return tile
		}
		accum += weight
	}
	return 0
}

func GenWorld(nLayers int, width int, height int, seed *WorldSeed) *World {
	world := World{}
	world.Width = width
	world.Height = height
	world.Layers = make(Layers, nLayers)
	layers := world.Layers
	for z := 0; z < nLayers; z++ {
		layers[z] = make(Layer, height)
		for y := 0; y < height; y++ {
			layers[z][y] = make(TileRow, width)
			for x := 0; x < width; x++ {
				tile := seed.NextTile()
				layers[z][y][x] = tile
			}
		}
	}
	return &world
}

func main() {
	var world *World
	width := flag.Int("width", 50, "width of world")
	height := flag.Int("height", 100, "height of world")
	grass := flag.Int("grass", 10, "weight of grass tiles")
	dirt := flag.Int("dirt", 16, "weight of dirt tiles")
	lava := flag.Int("lava", 1, "weight of lava tiles")
	water := flag.Int("water", 2, "weight of water tiles")
	generateImage := flag.Bool("image", false, "create an image")
	fileName := flag.String("image-name", "out", "weight of wate")

	flag.Parse()

	seed := WorldSeed{TILE_GRASS: *grass, TILE_DIRT: *dirt, TILE_LAVA: *lava, TILE_WATER: *water}
	world = GenWorld(1, *width, *height, &seed)

	if *generateImage == true {
		worldBigger := world.Scale(2);
		RenderMap(*fileName, worldBigger)
	} else {
		fmt.Print(world)
	}
}
