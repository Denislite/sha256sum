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
}

func main() {
	switch {
	case len(path) > 0:
		internal.Sha256Sum(path)
	default:
		log.Println(internal.ErrorOption)
	}
}
