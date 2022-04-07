package service

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sha256sum/internal/utils"
)

var (
	paths  = make(chan string)
	hashes = make(chan string)
)

func Initialize(dir, path, hashType string, help bool, ctx context.Context) {
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
		go Sha256sum(paths, hashes, hashType)
		go LookUpManager(dir, paths)
		PrintResult(hashes, ctx)

	case len(path) > 0:
		fmt.Println(FileHash(path, hashType))
	default:
		log.Println(utils.ErrorOption)
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
				log.Println(utils.ErrorDocs)
				return
			}
		})
	}
}

func CheckSignal(signals chan os.Signal) {
	signal.Notify(signals, os.Interrupt)
	go func() {
		for sig := range signals {
			log.Printf("request canceled by signal %d \n", sig)
			os.Exit(1)
		}
	}()
}
