package test

import (
	"sha256sum/internal"
	"sync"
	"testing"
)

const testPath = "your_path"

// speed up about 100x with Goroutines and WaitGroup using
func BenchmarkDefault(b *testing.B) {
	var wg sync.WaitGroup
	wg.Add(1)

	b.ResetTimer()
	go internal.SearchFiles(testPath, &wg)
	b.StopTimer()

	wg.Wait()
}
