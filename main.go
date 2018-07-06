package main

import (
	"log"
	"net/http"
	"os"

	"github.com/tusupov/gomandelbrot/handle"
)

var argv struct {
	listen 		string
}

func initArgv() {

	argv.listen = ":8080"

	listen := os.Getenv("LISTEN")
	if len(listen) > 0 {
		argv.listen = listen
	}

}

func main() {

	initArgv()

	// Обработчик алгоритма "Множество Мандельброта"
	http.HandleFunc("/mandelbrot/", handle.Mandelbrot)

	// Запускаем сервер
	log.Fatal(http.ListenAndServe(argv.listen, nil))

}
