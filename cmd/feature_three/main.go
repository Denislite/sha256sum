package main

import (
	"flag"
	"log"
	"sha256sum/internal"
)

var (
	path   string
	paths  chan string
	hashes chan string
)

func init() {
	flag.StringVar(&path, "o", "", "directory path")
	flag.Parse()
	paths = make(chan string)
	hashes = make(chan string)
}

func main() {
	switch {
	case len(path) > 0:

		go internal.Sha256sum(paths, hashes)
		go internal.LookUpManager(path, paths)
		internal.PrintResult(hashes)

	default:
		log.Println(internal.ErrorOption)
	}
}
