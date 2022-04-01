package internal

import (
	"regexp"
)

const (
	pattern = ".+\\/"
)

func checkFilePath(path, name string) string {
	matched, _ := regexp.MatchString(pattern, path)
	if matched != true {
		return path + "/" + name
	}
	return path + name
}
