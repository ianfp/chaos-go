package main

import (
	"github.com/ianfp/chaos/mandelbrot"
	"os"
)

func main() {
	query := ""
	if len(os.Args) > 1 {
		query = os.Args[1]
	}
	if query == "serve" {
		mandelbrot.Serve()
	} else {
		mandelbrot.Mandelbrot(mandelbrot.ParseCli(query))
	}
}

