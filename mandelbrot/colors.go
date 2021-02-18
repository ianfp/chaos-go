package mandelbrot

import "image/color"

func grayscale(iter int) color.Color {
	greyScale := 255 - uint8(255*iter/maxIter)
	return color.RGBA{R: greyScale, G: greyScale, B: greyScale, A: 255}
}

func fullColor(iter int) color.Color {
	if iter == maxIter {
		return color.Black // it's in the set
	}
	return color.RGBA{
		R: colorFor(iter, 0),
		G: colorFor(iter, maxIter / 3),
		B: colorFor(iter, 2 * maxIter / 3),
		A: 255,
	}
}

func colorFor(iter int, offset int) uint8 {
	x := iter + offset
	slope := 6 * 255 / maxIter
	yIntercept := -255
	if x > maxIter / 2 {
		slope = -slope
		yIntercept = 5 * 255
	}
	y := (slope * x) + yIntercept
	if y > 255 {
		y = 255
	} else if y < 0 {
		y = 0
	}
	return uint8(y)
}
