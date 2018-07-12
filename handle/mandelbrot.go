package handle

import (
	"net/http"
	"strconv"

	"github.com/tusupov/gomandelbrot/mandelbrot"
)

// Обработчик алгоритма "Множество Мандельброта"
func Mandelbrot(w http.ResponseWriter, r *http.Request) {

	// Параметры запросов
	queryParams := r.URL.Query()

	// Получаем нужные параметры
	x, _ := strconv.ParseFloat(queryParams.Get("x"), 10)
	y, _ := strconv.ParseFloat(queryParams.Get("y"), 10)
	zoom, _ := strconv.ParseUint(queryParams.Get("zoom"), 10, 0)
	rec := queryParams.Get("rec")

	// Создаем и получаем "Множество Мандельброта"
	// в виде байт массива
	buf := mandelbrot.
		New(rec).
		Move(x, y).
		Zoom(zoom).
		Draw()

	// Выводим на экран картинку
	writeImage(w, buf)

}

// Показать картинку
func writeImage(w http.ResponseWriter, buf []byte) {

	// Устанавливаем нужные параметры для header
	w.Header().Set("Content-Disposition", "inline; filename=\"mandelbrot.png\"")
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(buf)))

	w.WriteHeader(http.StatusOK)

	// Записываем картинку
	w.Write(buf)

}
