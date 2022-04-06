package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sha256sum/internal"
)

var (
	path   string
	help   bool
	paths  chan string
	hashes chan string
)

func init() {
	flag.StringVar(&path, "o", "", "/example/.../dir/ || you can check hash sum by dir/file path")
	flag.BoolVar(&help, "h", false, "|| read documentation")
	flag.Parse()

	paths = make(chan string)
	hashes = make(chan string)
}

func main() {
	switch {
	case help:
		createDocs()
		flag.Usage()

	case len(path) > 0:
		go internal.Sha256sum(paths, hashes)
		go internal.LookUpManager(path, paths)
		internal.PrintResult(hashes)

	default:
		log.Println(internal.ErrorOption)
	}
}

func createDocs() {
	flag.Usage = func() {
		_, err := fmt.Fprintln(os.Stderr, "Options for tool:")
		if err != nil {
			return
		}
		flag.VisitAll(func(f *flag.Flag) {
			_, err := fmt.Fprintf(os.Stderr, "-%v %v  \n", f.Name, f.Usage)
			if err != nil {
				return
			}
		})
	}
}
