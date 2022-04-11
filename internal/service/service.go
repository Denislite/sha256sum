package service

import "sha256sum/internal/repository"

type Hasher interface {
	FileHash(path, hashType string) (string, error)
	DirectoryHash()
}

type Service struct {
	Hasher
}

func NewService(repo repository.Repository) *Service {
	return &Service{Hasher: NewHasherService(repo)}
}
