package internal

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
)

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
