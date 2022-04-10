package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/vend/vend"
	"github.com/vend/vend/http"
	"github.com/vend/vend/inmemory"
)

func main() {
	port := flag.String("port", "8080", "port to http server, default 8080.")
	storage := inmemory.NewInMemory()
	seedMockData(storage)

	s := http.NewServer(*port, storage)

	fmt.Printf("Starting server on port %s.\n", *port)
	err := s.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func seedMockData(s vend.Storage) error {

	_, _ = s.CreateProduct("Golf Club", 100000)
	_, _ = s.CreateProduct("Nike Shorts", 1000)
	_, _ = s.CreateProduct("Toothbrush", 1000)

	return nil
}
