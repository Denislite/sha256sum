package test

import (
	"fmt"
	"sha256sum/internal"
	"testing"
)

const testPath = "/Users/denislogvinov/Downloads/smth/test"

func BenchmarkConcurrency(b *testing.B) {
	b.ResetTimer()

	paths := make(chan string)
	hashes := make(chan string)

	go internal.Sha256sum(paths, hashes)
	go internal.LookUpManager(testPath, paths)
	internal.PrintResult(hashes)

	b.StopTimer()
}

func BenchmarkDefault(b *testing.B) {
	b.ResetTimer()
	value, err := internal.DirectoryHash(testPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	for file, hash := range value {
		fmt.Printf("file %s || checksum: %s \n", file, hash)
	}
	b.StopTimer()
}
