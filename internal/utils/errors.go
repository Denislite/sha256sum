package utils

import "errors"

var (
	ErrorOption       = errors.New("error occurred with options, use -h to learn more")
	ErrorDocs         = errors.New("error occurred while creating documentation")
	ErrorDBConnection = errors.New("error occurred while creating connect to DB")
	ErrorConfig       = errors.New("error occurred while parsing config file")
)
