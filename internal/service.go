package internal

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

// FileHash - function to get file hash sum without concurrency
func FileHash(path string) (string, error) {
	file, err := os.Open(path)

	if err != nil {
		return "", ErrorWrongFile
	}

	defer file.Close()

	hash := sha256.New()
	_, err = io.Copy(hash, file)

	if err != nil {
		return "", ErrorHash
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// DirectoryHash - function to get files(from dir) hash sum without concurrency
func DirectoryHash(path string) (map[string]string, error) {
	filesHash := make(map[string]string)

	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return ErrorDirectoryRead
			}
			if info.IsDir() == false {
				value, err := FileHash(path)
				if err != nil {
					return err
				}
				filesHash[path] = value
			}
			return nil
		})
	if err != nil {
		return nil, ErrorDirectoryRead
	}

	return filesHash, nil
}

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
	}
}

func Hasher(wg *sync.WaitGroup, paths <-chan string, hashes chan<- string) {
	defer wg.Done()
	for path := range paths {
		file, err := os.Open(path)

		if err != nil {
			log.Println(ErrorWrongFile)
			return
		}

		defer file.Close()

		hash := sha256.New()
		_, err = io.Copy(hash, file)

		if err != nil {
			log.Println(ErrorHash)
			return
		}

		hashes <- fmt.Sprintf("file %s || checksum: %s", path, hex.EncodeToString(hash.Sum(nil)))
	}
}

func Sha256sum(paths, hashes chan string) {
	var wg sync.WaitGroup
	for worker := 1; worker <= (runtime.NumCPU() / 2); worker++ {
		wg.Add(1)
		go Hasher(&wg, paths, hashes)
	}
	defer close(hashes)
	wg.Wait()
}

func PrintResult(hashes chan string) {
	for {
		select {
		case hash, ok := <-hashes:
			if !ok {
				return
			}
			fmt.Println(hash)
		}
	}
}
