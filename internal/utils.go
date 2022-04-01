package internal

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
)

func FileHash(path string) string {
	file, err := os.Open(path)

	if err != nil {
		log.Println(fileError)
		return ""
	}

	defer file.Close()

	hash := sha256.New()
	_, err = io.Copy(hash, file)

	if err != nil {
		log.Println(hashError)
		return ""
	}

	return hex.EncodeToString(hash.Sum(nil))
}

func DirectoryHash(path string) string {
	fmt.Printf("dir %s \n", path)
	return ""
}
