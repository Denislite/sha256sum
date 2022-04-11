package service

import (
	"context"
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
		return "", err
	}

	return value, nil
}

func (s HasherService) DirectoryHash(ctx context.Context, path, hashType string) error {
	paths := make(chan string)
	hashes := make(chan string)

	go hash.Sha256sum(paths, hashes, hashType)
	go hash.LookUpManager(path, paths)
	hash.PrintResult(ctx, hashes)

	return nil
}
