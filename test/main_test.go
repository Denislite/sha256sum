package test

import (
	"fmt"
	"sha256sum/internal"
	"testing"
)

const testPath = "/Users/denislogvinov/Downloads/javafx-sdk-17.0.2"

//speed up about 3-3.5x with 3 workers
func BenchmarkConcurrency(b *testing.B) {
	b.ResetTimer()

	internal.Sha256Sum(testPath)

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
