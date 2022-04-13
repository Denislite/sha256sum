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

func (s HasherService) FileHash(path, hashType string) (*hash.FileInfo, error) {

	value, err := hash.FileHash(path, hashType)

	if err != nil {
		return nil, err
	}

	err = s.repo.SaveHash(*value)

	if err != nil {
		return nil, err
	}

	return value, nil
}

func (s HasherService) DirectoryHash(ctx context.Context, path, hashType string) error {
	paths := make(chan string)
	hashes := make(chan hash.FileInfo)

	go hash.Sha256sum(paths, hashes, hashType)
	go hash.LookUpManager(path, paths)
	value := hash.PrintResult(ctx, hashes)

	err := s.repo.SaveDirectoryHash(value)
	if err != nil {
		return err
	}

	return nil
}
