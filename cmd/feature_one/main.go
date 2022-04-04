package main

import (
	"flag"
	"fmt"
	"log"
	"sha256sum/internal"
)

var (
	file string
)

func init() {
	flag.StringVar(&file, "f", "", "file path")
	flag.Parse()
}

func main() {
	switch {
	case len(file) > 0:
		value, err := internal.FileHash(file)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Printf("file %s checksum: %s \n", file, value)
	default:
		log.Println("error based on command syntax")
	}
}
