package internal

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
)

func FileHash(path string) {
	file, err := os.Open(path)

	if err != nil {
		log.Println(fileError)
		return
	}

	defer file.Close()

	hash := sha256.New()
	_, err = io.Copy(hash, file)

	if err != nil {
		log.Println(hashError)
		return
	}

	fmt.Printf("file %s || checksum: %s \n", path, hex.EncodeToString(hash.Sum(nil)))
}

func DirectoryHash(path string) {
	dir, err := os.Open(path)

	if err != nil {
		log.Println(fileError)
		return
	}

	defer dir.Close()

	files, err := dir.ReadDir(0)

	if err != nil {
		log.Println(dirError)
		return
	}

	for _, v := range files {
		FileHash(checkFilePath(path, v.Name()))
	}
}
