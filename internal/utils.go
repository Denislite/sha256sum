package internal

import (
	"regexp"
)

const (
	pattern = ".+\\/"
)

func checkFilePath(path, name string) string {
	matched, _ := regexp.MatchString(pattern, path)
	if matched != false {
		return path + "/" + name
	}
	return path + name
}
