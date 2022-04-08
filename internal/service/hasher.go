package service

import "sha256sum/pkg/hash"

type HasherService struct {
}

func NewHasherService() *HasherService {
	return &HasherService{}
}

func (h HasherService) FileHash(path, hashType string) (string, error) {

	value := hash.FileHash(path, hashType)

	return value, nil
}

func (h HasherService) DirectoryHash() {
	//TODO implement me
	panic("implement me")
}
