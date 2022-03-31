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
		log.Fatalln(fileError)
	}

	defer file.Close()

	hash := sha256.New()
	_, err = io.Copy(hash, file)

	if err != nil {
		log.Fatalln(hashError)
	}

	fmt.Printf("file %s checksum: %s \n", path, hex.EncodeToString(hash.Sum(nil)))
}

func DirectoryHash(path string) {
	fmt.Printf("dir %s \n", path)
	return
}
