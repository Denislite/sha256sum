package hash

import (
	"context"
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sha256sum/internal/model"
	"sync"

	"sha256sum/internal/utils"
)

// FileHash - function to get file hash sum
func FileHash(path, hashType string) (*model.Hasher, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, utils.ErrorWrongFile
	}

	defer file.Close()

	data := model.Hasher{
		FileName: filepath.Base(path),
		FilePath: path,
		HashType: hashType,
	}

	switch hashType {
	case "md5":
		hash := md5.New()
		_, err = io.Copy(hash, file)
		data.HashValue = hash.Sum(nil)
	case "sha512":
		hash := sha512.New()
		_, err = io.Copy(hash, file)
		data.HashValue = hash.Sum(nil)
	default:
		hash := sha256.New()
		_, err = io.Copy(hash, file)
		data.HashValue = hash.Sum(nil)
		data.HashType = "sha256"
	}

	if err != nil {
		return nil, utils.ErrorHash
	}

	return &data, nil
}

// LookUpManager - function to get files path
func LookUpManager(inputPath string, paths chan string) {
	err := filepath.Walk(inputPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return utils.ErrorDirectoryRead
		}
		if !info.IsDir() {
			paths <- path
		}

		return nil
	})
	close(paths)

	if err != nil {
		log.Println(utils.ErrorDirectoryRead)
		return
	}
}

// Hasher - function to get all files hashes from directory
func Hasher(wg *sync.WaitGroup, paths <-chan string, hashes chan<- model.Hasher, hashType string) {
	defer wg.Done()
	for path := range paths {
		hash, err := FileHash(path, hashType)
		if err != nil {
			log.Println(err)
		}
		hashes <- *hash
	}
}

// Sha256sum - main function which init our workers pool
func Sha256sum(paths chan string, hashes chan model.Hasher, hashType string) {
	var wg sync.WaitGroup
	for worker := 1; worker <= runtime.NumCPU(); worker++ {
		wg.Add(1)
		go Hasher(&wg, paths, hashes, hashType)
	}
	defer close(hashes)
	wg.Wait()
}

// PrintResult - output function
func PrintResult(ctx context.Context, hashes chan model.Hasher) []model.Hasher {
	var result []model.Hasher
	for {
		select {
		case hash, ok := <-hashes:
			if !ok {
				return result
			}
			result = append(result, hash)
			fmt.Printf("%x %s \n", hash.HashValue, hash.FileName)
		case <-ctx.Done():
			log.Println("request canceled by context")
			os.Exit(1)
			return nil
		}
	}
	return result
}
