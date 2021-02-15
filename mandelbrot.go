package main

import (
	"fmt"
	"github.com/fogleman/gg"
	"image/color"
	"math"
)

// Image size in pixels
const imageSize = 1000

// Number of iterations before concluding that the point
// is probably in the Mandelbrot set.
const maxIter = 200

// Describes the section of the complex plane that we're viewing.
type viewport struct {
	topLeft complex128
	bottomRight complex128
}

func (vp viewport) width() float64 {
	return real(vp.bottomRight) - real(vp.topLeft)
}

func (vp viewport) height() float64 {
	return imag(vp.topLeft) - imag(vp.bottomRight)
}

func (vp viewport) pointAt(fromLeft int, fromTop int) complex128 {
	xDelta := float64(fromLeft) / imageSize * vp.width()
	yDelta := float64(fromTop) / imageSize * vp.height()
	realPart := real(vp.topLeft) + xDelta
	imagPart := imag(vp.topLeft) - yDelta
	return complex(realPart, imagPart)
}

func mandelbrot() {
	vp := viewport{-1 + 1i, 1 - 1i}
	dc := gg.NewContext(gridSize, gridSize)
	for fromLeft := 0; fromLeft < imageSize; fromLeft++ {
		for fromTop := 0; fromTop < imageSize; fromTop++ {
			value := vp.pointAt(fromLeft, fromTop)
			dc.SetColor(colorOfIter(stepsBeforeDiverge(value)))
			dc.SetPixel(fromLeft, fromTop)
		}
	}
	err := dc.SavePNG("mandelbrot.png")
	if err != nil {
		fmt.Printf("ERROR: failed to save PNG: %v", err)
	}
}

func colorOfIter(iter int) color.Color {
	greyScale := 255 - uint8(255 * iter / maxIter)
	return color.RGBA{R: greyScale, G: greyScale, B: greyScale, A: 255}
}

func stepsBeforeDiverge(value complex128) int {
	current := 0 + 0i
	for iter := 0; iter < maxIter; iter++ {
		if diverges(current) {
			return iter
		}
		current = (current * current) + value
	}
	return maxIter
}

func diverges(current complex128) bool {
	return math.Abs(real(current)) > 2 || math.Abs(imag(current)) > 2
}
