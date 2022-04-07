package main

import (
	"context"
	"flag"
	"os"
	"sha256sum/internal/service"
)

var (
	dir      string
	path     string
	hashType string
	help     bool
	ctx      context.Context
	signals  chan os.Signal
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
	service.CheckSignal(signals)
	service.Initialize(dir, path, hashType, help, ctx)
}
