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
	deleted  string
	help     bool
	ctx      context.Context
	signals  chan os.Signal
)

func init() {
	flag.StringVar(&dir, "d", "", "/example/.../dir/ || you can check hashsum sum by dir path")
	flag.StringVar(&path, "f", "", "/example/.../text.txt || you can check hashsum sum by file path")
	flag.StringVar(&hashType, "a", "sha256", "available: md5, sha512 || default: sha256")
	flag.StringVar(&check, "check", "", "check old hash in db with new one")
	flag.StringVar(&deleted, "deleted", "", "check deleted files if dir")
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
		return
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
		resultHash, err := s.Hasher.CompareHash(ctx, check, hashType)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println("files hash changes (old/new):")
		for _, hash := range resultHash {
			fmt.Printf("%s %s || %s \n",
				hash.FileName, hash.OldHash, hash.NewHash)
		}

	case len(deleted) > 0:
		resultFiles, err := s.Hasher.CheckDeleted(ctx, deleted, hashType)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println("deleted files:")
		for _, files := range resultFiles {
			fmt.Printf("%s %s \n", files.FilePath, files.OldHash)
		}
	default:
		log.Println(utils.ErrorOption)
	}
}
