package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sha256sum/internal"
)

var (
	dir      string
	path     string
	hashType string
	help     bool
	paths    chan string
	hashes   chan string
	signals  chan os.Signal
	ctx      context.Context
)

func init() {
	flag.StringVar(&dir, "d", "", "/example/.../dir/ || you can check hash sum by dir path")
	flag.StringVar(&path, "f", "", "/example/.../text.txt || you can check hash sum by file path")
	flag.StringVar(&hashType, "a", "", "available: md5, sha512 || default: sha256")
	flag.BoolVar(&help, "h", false, "|| you can read options")
	flag.Parse()

	paths = make(chan string)
	hashes = make(chan string)
	signals = make(chan os.Signal, 1)
	ctx = context.Background()
}

func main() {

	signal.Notify(signals, os.Interrupt)
	go func() {
		for sig := range signals {
			log.Printf("request canceled by signal %d \n", sig)
			os.Exit(1)
		}
	}()

	ctx, cancel := context.WithCancel(ctx)
	go func() {
		fmt.Scanln()
		cancel()
	}()

	switch {
	case help:
		createDocs()
		flag.Usage()

	case len(dir) > 0:
		go internal.Sha256sum(paths, hashes, hashType)
		go internal.LookUpManager(dir, paths)
		internal.PrintResult(hashes, ctx)

	case len(path) > 0:
		fmt.Println(internal.FileHash(path, hashType))
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
