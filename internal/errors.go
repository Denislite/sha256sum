package internal

import "errors"

var (
	fileError = errors.New("wrong file/dir path")
	hashError = errors.New("error occurred while taking a hash from file")
)
