package mandelbrot

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

func handler(w http.ResponseWriter, r *http.Request) {
	logInfo.Printf("Received request %v", r.URL)
	vp, parseFailed := parseQueryValues(r.URL.Query())
	if parseFailed != nil {
		badRequest(w, parseFailed)
		return
	}
	image := Mandelbrot(vp)
	w.Header().Add("Content-Type", "image/png")
	encodeFailed := image.EncodePNG(w)
	if encodeFailed != nil {
		logError.Printf("Encoding PNG failed: %s", encodeFailed)
	}
}

func parseQueryValues(query url.Values) (*viewport, error) {
	center := 0 + 0i
	var badCenter error = nil
	if centerString := query.Get("center"); centerString != "" {
		center, badCenter = strconv.ParseComplex(centerString, 64)
		if badCenter != nil {
			return nil, badCenter
		}
	}

	width := 2.0
	var badWidth error = nil
	if widthString := query.Get("width"); widthString != "" {
		width, badWidth = strconv.ParseFloat(widthString, 64)
		if badWidth != nil {
			return nil, badWidth
		}
	}

	return &viewport{
		topLeft: center + complex(-width, width),
		bottomRight: center + complex(width, -width),
	}, nil
}

func badRequest(w http.ResponseWriter, err error) {
	logError.Printf("Bad request: %v", err)
	w.WriteHeader(400)
	_, fatal := fmt.Fprint(w, err.Error())
	if fatal != nil {
		w.WriteHeader(500)
		log.Fatal(fatal)
	}
}

func Serve() {
	addr := ":9000"
	logInfo.Printf("Listening on %s", addr)
	http.HandleFunc("/", handler)
	logError.Fatal(http.ListenAndServe(addr, nil))
}

