package mandelbrot

import (
	"fmt"
	"github.com/fogleman/gg"
	"image/color"
	"log"
	"math"
	"net/url"
	"os"
	"strconv"
	"time"
)

// Image size in pixels
const imageSize = 1000

// Number of iterations before concluding that the point
// is probably in the Mandelbrot set.
const maxIter = 1000

const fontHeight = 10

// Describes the section of the complex plane that we're viewing.
type viewport struct {
	topLeft     complex128
	bottomRight complex128
}

var (
	logInfo = log.New(os.Stderr, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	logError = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
)

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

func Mandelbrot(vp *viewport) *gg.Context {
	startTime := time.Now()
	dc := gg.NewContext(imageSize, imageSize)
	for fromLeft := 0; fromLeft < imageSize; fromLeft++ {
		for fromTop := 0; fromTop < imageSize; fromTop++ {
			value := vp.pointAt(fromLeft, fromTop)
			dc.SetColor(fullColor(stepsBeforeDiverge(value)))
			dc.SetPixel(fromLeft, fromTop)
		}
	}
	elapsed := time.Now().Sub(startTime)
	logInfo.Printf("Calculated in %v", elapsed)

	displayCoords(vp, dc)
	err := dc.SavePNG("mandelbrot.png")
	if err != nil {
		logError.Printf("Failed to save PNG: %v", err)
	}
	return dc
}

func displayCoords(vp *viewport, dc *gg.Context) {
	dc.SetColor(color.RGBA{R: 255, A: 255})
	dc.DrawString(formatComplex128(vp.topLeft), 1, fontHeight)
	dc.DrawString(fmt.Sprintf("width=%e", vp.width()), imageSize-150, imageSize-fontHeight)
}

func formatComplex128(value complex128) string {
	return strconv.FormatComplex(value, 'g', 4, 64)
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

func defaultViewport() *viewport {
	return &viewport{-2 + 2i, 2 - 2i}
}

// CLI processing

func ParseCli(query string) *viewport {
	if query == "" {
		logInfo.Print("Using default viewport")
		return defaultViewport()
	}
	values, badQuery := url.ParseQuery(query)
	if badQuery != nil {
		logError.Printf("Unable to parse query %s: %s", query, badQuery)
		return defaultViewport()
	}
	fromQuery, badValues := parseQueryValues(values)
	if badValues != nil {
		logError.Printf("Bad query values: %s", badValues)
		return defaultViewport()
	}
	return fromQuery
}
