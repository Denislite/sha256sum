package utils

import "errors"

var (
	ErrorOption        = errors.New("error occurred with options, use -h to learn more")
	ErrorDocs          = errors.New("error occurred while creating documentation")
	ErrorDBConnection  = errors.New("error occurred while creating connect to DB")
	ErrorWrongFile     = errors.New("error occurred with file, please check your file path")
	ErrorHash          = errors.New("error occurred while taking a hash sum from file")
	ErrorDirectoryRead = errors.New("error occurred while reading files from directory")
)
