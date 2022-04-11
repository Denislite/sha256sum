package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"sha256sum/internal/utils"
	"sha256sum/pkg/hash"
)

var (
	dir      string
	path     string
	hashType string
	help     bool
	ctx      context.Context
	signals  chan os.Signal
	paths    = make(chan string)
	hashes   = make(chan string)
)

func init() {
	flag.StringVar(&dir, "d", "", "/example/.../dir/ || you can check hash sum by dir path")
	flag.StringVar(&path, "f", "", "/example/.../text.txt || you can check hash sum by file path")
	flag.StringVar(&hashType, "a", "", "available: md5, sha512 || default: sha256")
	flag.BoolVar(&help, "h", false, "|| you can read options")
	flag.Parse()

	signals = make(chan os.Signal, 1)
	ctx = context.Background()
}

func main() {
	utils.CheckSignal(signals)

	ctx, cancel := context.WithCancel(ctx)

	go func() {
		_, err := fmt.Scanln()
		if err != nil {
			return
		}
		cancel()
	}()

	switch {
	case help:
		utils.CreateDocs()
		flag.Usage()

	case len(dir) > 0:
		go hash.Sha256sum(paths, hashes, hashType)
		go hash.LookUpManager(dir, paths)
		hash.PrintResult(ctx, hashes)

	case len(path) > 0:
		hash, err := hash.FileHash(path, hashType)
		if err != nil {
			return
		}
		fmt.Println(hash)

	default:
		log.Println(utils.ErrorOption)
	}
}
