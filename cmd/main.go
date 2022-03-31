package main

import (
	"flag"
	"log"
	"sha256sum/internal"
)

func main() {
	var file string
	var dir string

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
