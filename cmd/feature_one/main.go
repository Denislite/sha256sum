package main

import (
	"flag"
	"fmt"
	"log"
	"sha256sum/internal"
)

var (
	file string
	dir  string
)

func init() {
	flag.StringVar(&file, "f", "", "file path")
	flag.Parse()
}

func main() {
	switch {
	case len(dir) > 0:
		value, err := internal.FileHash(dir)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Printf("file %s checksum: %s \n", dir, value)
	default:
		log.Println("error based on command syntax")
	}
}
