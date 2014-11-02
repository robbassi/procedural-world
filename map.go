package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
)

const (
	TILE_WIDTH int = 16

	O_TOP_LEFT     = 0
	O_TOP          = 1
	O_TOP_RIGHT    = 2
	O_RIGHT        = 3
	O_BOTTOM       = 4
	O_BOTTOM_RIGHT = 5
	O_BOTTOM_LEFT  = 6
	O_LEFT         = 7

	I_TOP_LEFT     = 8
	I_TOP_RIGHT    = 9
	I_BOTTOM_LEFT  = 10
	I_BOTTOM_RIGHT = 11
)

func RenderMap(mapName string, world *World) {
	coreTiles := map[Tile]image.Image{
		TILE_GRASS: loadImg("grass"),
		TILE_DIRT:  loadImg("dirt"),
		TILE_LAVA:  loadImg("lava"),
		TILE_WATER: loadImg("water"),
	}
	// corners
	fluidTiles := map[Tile][]image.Image{
		TILE_WATER: []image.Image{
			loadImg("waterInTL"),
			loadImg("waterInT"),
			loadImg("waterInTR"),
			loadImg("waterInR"),
			loadImg("waterInB"),
			loadImg("waterInBR"),
			loadImg("waterInBL"),
			loadImg("waterInL"),
			loadImg("waterOutTL"),
			loadImg("waterOutTR"),
			loadImg("waterOutBL"),
			loadImg("waterOutBR"),
		},
		TILE_LAVA: []image.Image{
			loadImg("lavaInTL"),
			loadImg("lavaInT"),
			loadImg("lavaInTR"),
			loadImg("lavaInR"),
			loadImg("lavaInB"),
			loadImg("lavaInBR"),
			loadImg("lavaInBL"),
			loadImg("lavaInL"),
			loadImg("lavaOutTL"),
			loadImg("lavaOutTR"),
			loadImg("lavaOutBL"),
			loadImg("lavaOutBR"),
		},
	}

	for _, layer := range world.Layers {
		layerImg := image.NewRGBA(image.Rect(0, 0, world.Width*TILE_WIDTH, world.Height*TILE_WIDTH))
		for y, row := range layer {
			for x, tile := range row {
				img := coreTiles[tile]
				rect := image.Rect(x*TILE_WIDTH, y*TILE_WIDTH, x*TILE_WIDTH+TILE_WIDTH, y*TILE_WIDTH+TILE_WIDTH)

				if layer.IsTile(x, y, TILE_LAND) {
					draw.Draw(layerImg, rect, img, image.Pt(0, 0), draw.Over)
				} else {
					var landTile Tile

					// Surrounded by fluid?
					if landTile = layer.ClosestLand(x, y); landTile == TILE_EMPTY {
						draw.Draw(layerImg, rect, img, image.Pt(0, 0), draw.Over)
					} else {

						// draw land tile before mask
						draw.Draw(layerImg, rect, coreTiles[landTile], image.Pt(0, 0), draw.Over)

						// Surrounded by land?
						if layer.IsTile(x-1, y, TILE_LAND) && layer.IsTile(x, y-1, TILE_LAND) {
							// top left
							draw.Draw(layerImg, rect, fluidTiles[tile][O_TOP_LEFT], image.Pt(0, 0), draw.Over)
						} else if layer.IsTile(x+1, y, TILE_LAND) && layer.IsTile(x, y-1, TILE_LAND) {
							// top right
							draw.Draw(layerImg, rect, fluidTiles[tile][O_TOP_RIGHT], image.Pt(0, 0), draw.Over)
						} else if layer.IsTile(x, y+1, TILE_LAND) && layer.IsTile(x-1, y, TILE_LAND) {
							// bottom left
							draw.Draw(layerImg, rect, fluidTiles[tile][O_BOTTOM_LEFT], image.Pt(0, 0), draw.Over)
						} else if layer.IsTile(x, y+1, TILE_LAND) && layer.IsTile(x+1, y, TILE_LAND) {
							// bottom right
							draw.Draw(layerImg, rect, fluidTiles[tile][O_BOTTOM_RIGHT], image.Pt(0, 0), draw.Over)
						} else if layer.IsTile(x, y-1, TILE_LAND) {
							// top
							draw.Draw(layerImg, rect, fluidTiles[tile][O_TOP], image.Pt(0, 0), draw.Over)
						} else if layer.IsTile(x+1, y, TILE_LAND) {
							// right
							draw.Draw(layerImg, rect, fluidTiles[tile][O_RIGHT], image.Pt(0, 0), draw.Over)
						} else if layer.IsTile(x, y+1, TILE_LAND) {
							// bottom
							draw.Draw(layerImg, rect, fluidTiles[tile][O_BOTTOM], image.Pt(0, 0), draw.Over)
						} else if layer.IsTile(x-1, y, TILE_LAND) {
							// left
							draw.Draw(layerImg, rect, fluidTiles[tile][O_LEFT], image.Pt(0, 0), draw.Over)
						} else {
							// Surrounding some piece of land?
							if layer.IsTile(x-1, y-1, TILE_LAND) {
								// top left
								draw.Draw(layerImg, rect, fluidTiles[tile][I_TOP_LEFT], image.Pt(0, 0), draw.Over)
							} else if layer.IsTile(x+1, y-1, TILE_LAND) {
								// top right
								draw.Draw(layerImg, rect, fluidTiles[tile][I_TOP_RIGHT], image.Pt(0, 0), draw.Over)
							} else if layer.IsTile(x-1, y+1, TILE_LAND) {
								// bottom left
								draw.Draw(layerImg, rect, fluidTiles[tile][I_BOTTOM_LEFT], image.Pt(0, 0), draw.Over)
							} else if layer.IsTile(x+1, y+1, TILE_LAND) {
								// bottom right
								draw.Draw(layerImg, rect, fluidTiles[tile][I_BOTTOM_RIGHT], image.Pt(0, 0), draw.Over)
							}
						}
					}
				}
			}
		}
		fmt.Println("Writing file...")
		file, err := os.Create(mapName + ".png")
		if err != nil {
			fmt.Println(err)
		}
		png.Encode(file, layerImg)
		file.Close()
	}
}

func loadImg(name string) image.Image {
	f, err := os.Open("./img/" + name + ".png")
	if err != nil {
		fmt.Println(err)
	}

	img, _, _ := image.Decode(f)
	return img
}
