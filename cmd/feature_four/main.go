package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"sha256sum/internal/configs"
	"sha256sum/internal/repository"
	"sha256sum/internal/service"
	"sha256sum/internal/utils"
)

var (
	dir      string
	path     string
	hashType string
	check    string
	help     bool
	ctx      context.Context
	signals  chan os.Signal
)

func init() {
	flag.StringVar(&dir, "d", "", "/example/.../dir/ || you can check hashsum sum by dir path")
	flag.StringVar(&path, "f", "", "/example/.../text.txt || you can check hashsum sum by file path")
	flag.StringVar(&hashType, "a", "sha256", "available: md5, sha512 || default: sha256")
	flag.StringVar(&check, "check", "", "check old hash in db with new one")
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

	cfg, err := configs.ParseConfigFile("./internal/configs")
	if err != nil {
		log.Println(utils.ErrorConfig)
	}

	db, err := repository.NewPostgresDB(cfg)
	if err != nil {
		log.Println(utils.ErrorDBConnection)
	}

	r := repository.NewRepository(db)

	s := service.NewService(r)

	switch {
	case help:
		utils.CreateDocs()
		flag.Usage()

	case len(dir) > 0:
		result, err := s.Hasher.DirectoryHash(ctx, dir, hashType)
		if err != nil {
			log.Println(err)
		}
		for _, hash := range result {
			fmt.Printf("%s %s \n", hash.HashValue, hash.FileName)
		}

	case len(path) > 0:
		hash, err := s.Hasher.FileHash(path, hashType)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Printf("%s %s \n", hash.HashValue, hash.FileName)

	case len(check) > 0:
		resultHash, resultDeleted, err := s.Hasher.CompareHash(ctx, check, hashType)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println("files hash changes:")
		for _, hash := range resultHash {
			fmt.Printf("%s %s || %s \n",
				hash.FileName, hash.OldHash, hash.NewHash)
		}
		fmt.Println("deleted files:")
		for _, del := range resultDeleted {
			fmt.Printf("%s %s \n",
				del.FileName, del.OldHash)
		}
	default:
		log.Println(utils.ErrorOption)
	}
}
