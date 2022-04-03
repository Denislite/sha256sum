package internal

import "errors"

var (
	ErrorWrongFile     = errors.New("wrong file path")
	ErrorHash          = errors.New("error occurred while taking a hash from file")
	ErrorDirectoryRead = errors.New("error occurred while reading files from directory")
	ErrorOption        = errors.New("no option, use -f your_path")
)
