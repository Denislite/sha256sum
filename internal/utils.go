package internal

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"os"
)

func FileHash(path string) (string, error) {
	file, err := os.Open(path)

	if err != nil {
		log.Println(fileError)
		return "", fileError
	}

	defer file.Close()

	hash := sha256.New()
	_, err = io.Copy(hash, file)

	if err != nil {
		log.Println(hashError)
		return "", hashError
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
