package main

import (
	"flag"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/tusupov/gomandelbrot/handle"
	"github.com/tusupov/gomandelbrot/mandelbrot"
)

func main() {

	defer func() {
		if e := recover(); e != nil {
			log.Printf("main() %s: \n%s", e, debug.Stack())
		}
	}()

	addr := flag.String("h", ":8080", "Host address")
	workers := flag.Uint("w", 2, "Workers count")
	config := flag.String("c", "sizes.json", "Sizes config file path")
	flag.Parse()

	// Set workers count
	err := mandelbrot.SetWorkers(*workers)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Workers: %d\n", *workers)

	// Load size
	log.Printf("Size config file path: %s\n", *config)
	n, err := mandelbrot.LoadSize(*config)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Find size count: %d\n", n)

	// Mandelbrot handle function
	http.HandleFunc("/mandelbrot/", handle.Mandelbrot)

	// Start server
	log.Printf("Listening [%s] ...\n", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))

}
