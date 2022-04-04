package main

import (
	"flag"
	"log"
	"sha256sum/internal"
	"sync"
)

var (
	path string
	wg   sync.WaitGroup
)

func init() {
	flag.StringVar(&path, "o", "", "directory path")
	flag.Parse()
	wg.Add(1)
}

func main() {
	switch {
	case len(path) > 0:
		go internal.SearchFiles(path, &wg)
		wg.Wait()
	default:
		log.Println(internal.ErrorOption)
	}
}
