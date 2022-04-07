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
	path    string
	help    bool
	paths   chan string
	hashes  chan string
	signals chan os.Signal
	ctx     context.Context
)

func init() {
	flag.StringVar(&path, "o", "", "/example/.../dir/ || you can check hash sum by dir/file path")
	flag.BoolVar(&help, "h", false, "|| read documentation")
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

	case len(path) > 0:
		go internal.Sha256sum(paths, hashes)
		go internal.LookUpManager(path, paths)
		internal.PrintResultCtx(hashes, ctx)

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
