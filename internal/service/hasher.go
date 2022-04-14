package service

import (
	"context"
	"sha256sum/internal/repository"
	"sha256sum/pkg/hashsum"
)

type HasherService struct {
	repo repository.Repository
}

func NewHasherService(repo repository.Repository) *HasherService {
	return &HasherService{repo: repo}
}

func (s HasherService) FileHash(path, hashType string) (*hashsum.FileInfo, error) {

	value, err := hashsum.FileHash(path, hashType)

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
	hashes := make(chan hashsum.FileInfo)

	go hashsum.Sha256sum(paths, hashes, hashType)
	go hashsum.LookUpManager(path, paths)
	value := hashsum.PrintResult(ctx, hashes)

	err := s.repo.SaveDirectoryHash(value)
	if err != nil {
		return err
	}

	return nil
}
