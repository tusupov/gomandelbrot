package handle

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/tusupov/gomandelbrot/mandelbrot"
)

var (
	errQueryXEmpty    = "The required `x` parameter are not provided."
	errQueryYEmpty    = "The required `y` parameter are not provided."
	errQueryZoomEmpty = "The required `zoom` parameter are not provided."
)

// Mandelbrot handle function
func Mandelbrot(w http.ResponseWriter, r *http.Request) {

	// Query params
	queryParams := r.URL.Query()

	// Get required parameters
	xParam := queryParams.Get("x")
	yParam := queryParams.Get("y")
	zoomParam := queryParams.Get("zoom")

	if len(xParam) == 0 {
		http.Error(w, errQueryXEmpty, http.StatusBadRequest)
		return
	}
	x, err := strconv.ParseFloat(xParam, 10)
	if err != nil {
		http.Error(w, fmt.Sprintf("x: %s", err), http.StatusBadRequest)
		return
	}

	if len(yParam) == 0 {
		http.Error(w, errQueryXEmpty, http.StatusBadRequest)
		return
	}
	y, err := strconv.ParseFloat(yParam, 10)
	if err != nil {
		http.Error(w, fmt.Sprintf("y: %s", err), http.StatusBadRequest)
		return
	}

	if len(zoomParam) == 0 {
		http.Error(w, errQueryZoomEmpty, http.StatusBadRequest)
		return
	}

	zoom, err := strconv.ParseUint(zoomParam, 10, 0)
	if err != nil {
		http.Error(w, fmt.Sprintf("zoom: %s", err), http.StatusBadRequest)
		return
	}

	rec := queryParams.Get("rec")

	// Create and Get image
	m, err := mandelbrot.New(rec)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	buf, err := m.Move(x, y).Zoom(zoom).Draw()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write Image
	writeImage(w, buf)

}

// Write image
func writeImage(w http.ResponseWriter, buf []byte) {

	// Set all needs header configs
	w.Header().Set("Content-Disposition", "inline; filename=\"mandelbrot.png\"")
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(buf)))

	w.WriteHeader(http.StatusOK)

	// Write image
	if _, err := w.Write(buf); err != nil {
		log.Println(err)
	}

}
