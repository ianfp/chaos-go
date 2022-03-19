package game

import (
	"fmt"
	"github.com/fogleman/gg"
	"image/color"
	"math"
	"math/rand"
)

type vector struct {
	x float64
	y float64
}

const gridSize = 1000
const numIter = 500000
const printEvery = numIter / 10

func game() {
	points := make([]vector, 0)
	next := vector{0, 0}

	for iter := 0; iter < numIter; iter++ {
		points = append(points, next)
		//fmt.Printf("grid %d, %d is %d\n", next.x, next.y, grid[next.x][next.y])
		next = getNext(next)
		//fmt.Printf("next is %v\n", next)
		if iter%printEvery == 0 {
			drawPicture(points, iter)
		}
	}
}

func makeGrid(size int) [][]uint {
	grid := make([][]uint, size)
	for x := range grid {
		grid[x] = make([]uint, size)
	}
	return grid
}

func getNext(next vector) vector {
	if flipCoin() {
		return heads(next)
	} else {
		return tails(next)
	}
}

func flipCoin() bool {
	return rand.Int()%2 == 0
}

func drawPicture(points []vector, iter int) {
	dc := gg.NewContext(gridSize, gridSize)
	grid := makeGrid(gridSize)
	for _, vec := range points {
		x, y := toCoords(vec)
		grid[x][y]++
	}

	for x, row := range grid {
		for y := range row {
			count := grid[x][y]
			dc.SetColor(colorOf(count))
			dc.SetPixel(x, y)
		}
	}
	err := dc.SavePNG(fmt.Sprintf("game%v.png", iter))
	if err != nil {
		fmt.Printf("ERROR: failed to save PNG: %v", err)
	}
}

func toCoords(vec vector) (int, int) {
	return mod(vec.x), mod(vec.y)
}

func mod(val float64) int {
	result := int(val) % gridSize
	for result < 0 {
		result += gridSize
	}
	return result
}

func colorOf(count uint) color.Color {
	count = count * 50
	if count > 255 {
		count = 255
	}
	value := uint8(count)
	return color.RGBA{R: value, G: value, B: value, A: 255}
}

func heads(prev vector) vector {
	return prev.turnAndMove(0.8, 100)
}

func tails(prev vector) vector {
	return vector{prev.x + 300, prev.y}
}

func (prev vector) turnAndMove(radians float64, magnitude float64) vector {
	return vector{prev.x + (math.Cos(radians) * magnitude), prev.y + (math.Sin(radians) * magnitude)}
}
