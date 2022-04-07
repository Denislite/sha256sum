package internal

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
	"sync"
)

// FileHash - function to get file hash sum
func FileHash(path, hashType string) string {
	file, err := os.Open(path)

	if err != nil {
		log.Println(ErrorWrongFile)
		return ""
	}

	defer file.Close()

	var value interface{}

	switch hashType {
	case "md5":
		hash := md5.New()
		_, err = io.Copy(hash, file)
		value = hash.Sum(nil)
	case "512":
		hash := sha512.New()
		_, err = io.Copy(hash, file)
		value = hash.Sum(nil)
	default:
		hash := sha256.New()
		_, err = io.Copy(hash, file)
		value = hash.Sum(nil)
	}

	if err != nil {
		log.Println(ErrorHash)
		return ""
	}

	return fmt.Sprintf("file %s || checksum: %x", path, value)
}

// LookUpManager - function to get files path
func LookUpManager(inputPath string, paths chan string) {
	err := filepath.Walk(inputPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return ErrorDirectoryRead
		}
		if !info.IsDir() {
			paths <- path
		}

		return nil
	})
	close(paths)

	if err != nil {
		log.Println(err)
		return
	}
}

// Hasher - function to get all files hashes from directory
func Hasher(wg *sync.WaitGroup, paths <-chan string, hashes chan<- string, hashType string) {
	defer wg.Done()
	for path := range paths {
		hashes <- FileHash(path, hashType)
	}
}

// Sha256sum - main function which init our workers pool
func Sha256sum(paths, hashes chan string, hashType string) {
	var wg sync.WaitGroup
	for worker := 1; worker <= (runtime.NumCPU() / 2); worker++ {
		wg.Add(1)
		go Hasher(&wg, paths, hashes, hashType)
	}
	defer close(hashes)
	wg.Wait()
}

// PrintResult - output function
func PrintResult(hashes chan string, ctx context.Context) {
	for {
		select {
		case hash, ok := <-hashes:
			if !ok {
				return
			}
			fmt.Println(hash)
		case <-ctx.Done():
			log.Println("request canceled by context")
			os.Exit(1)
			return
		}
	}
}
