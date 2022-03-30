package internal

import "fmt"

func FileHash(path string) {
	fmt.Printf("file %s", path)
	return
}

func DirectoryHash(path string) {
	fmt.Printf("dir %s \n", path)
	return
}
