package hashsum

import "io"

type HashSum interface {
	Hash(file io.Reader) (string, error)
}

func New(algo string) (hasher HashSum, hashType string, err error) {
	switch algo {
	case "md5", "MD5":
		hasher = NewMD5()
		hashType = algo
	case "sha512", "SHA512":
		hasher = NewSha512()
		hashType = algo
	default:
		hasher = NewSha256()
		hashType = "sha256"
	}
	return
}
