package internal

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func DirectoryHashGo(path string) chan string {
	channel := make(chan string)

	go FileHashGo(channel)
	go func() {
		filepath.Walk(path,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() {
					channel <- path
				}
				return nil
			})
		defer close(channel)
	}()

	return channel
}

func FileHashGo(c chan string) {
	path := <-c

	file, err := os.Open(path)

	if err != nil {

	}

	defer file.Close()

	hash := sha256.New()
	_, err = io.Copy(hash, file)

	if err != nil {

	}

	fmt.Println(hex.EncodeToString(hash.Sum(nil)))
}