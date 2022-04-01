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
	flag.StringVar(&dir, "d", "", "directory path")
	flag.Parse()
}

func main() {
	switch {
	case len(dir) > 0:
		value, err := internal.DirectoryHash(dir)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Directory files:")
		for file, hash := range value {
			fmt.Printf("file %s || checksum: %s \n", file, hash)
		}
	case len(file) > 0:
		value, err := internal.FileHash(file)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("file %s || checksum: %s \n", file, value)
	default:
		log.Println("error based on command syntax")
	}
}
