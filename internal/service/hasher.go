package service

import (
	"log"
	"sha256sum/internal/repository"
	"sha256sum/pkg/hash"
)

type HasherService struct {
	repo repository.Repository
}

func NewHasherService(repo repository.Repository) *HasherService {
	return &HasherService{repo: repo}
}

func (s HasherService) FileHash(path, hashType string) (string, error) {

	value, err := hash.FileHash(path, hashType)

	if err != nil {
		log.Println(err)
		return "", err
	}

	return value, nil
}

func (s HasherService) DirectoryHash() {
	//TODO implement me
	panic("implement me")
}
