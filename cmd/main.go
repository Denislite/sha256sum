package main

import (
	"flag"
	"log"
	"sha256sum/internal"
)

var (
	file string
	dir  string
)

func main() {
	flag.StringVar(&file, "f", "", "file path")
	flag.StringVar(&dir, "d", "", "directory path")
	flag.Parse()

	switch {
	case len(dir) > 0:
		internal.DirectoryHash(dir)
	case len(file) > 0:
		internal.FileHash(file)
	default:
		log.Println("error based on command syntax")
	}
}
