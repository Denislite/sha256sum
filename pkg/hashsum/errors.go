package hashsum

import "errors"

var (
	ErrorWrongFile     = errors.New("error occurred with file, please check your file path")
	ErrorHash          = errors.New("error occurred while taking a hashsum from file")
	ErrorDirectoryRead = errors.New("error occurred while reading files from directory")
)
