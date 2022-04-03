package internal

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
)

func Initialize(path string, wg *sync.WaitGroup) {
	defer wg.Done()

	wg.Add(1)

	go SearchFiles(path, wg)
}

func TakeFileHash(path string, wg *sync.WaitGroup) {
	defer wg.Done()

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

	fmt.Printf("file %s || checksum: %s \n", path, hex.EncodeToString(hash.Sum(nil)))
}

func SearchFiles(commonPath string, wg *sync.WaitGroup) {
	defer wg.Done()

	err := filepath.Walk(commonPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return ErrorDirectoryRead
		}

		if !info.IsDir() {
			wg.Add(1)
			go TakeFileHash(path, wg)
		}
		return nil
	})

	if err != nil {
		log.Println(err)
		return
	}
}
