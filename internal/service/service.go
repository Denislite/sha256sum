package service

import (
	"context"
	"sha256sum/internal/repository"
	"sha256sum/pkg/hash"
)

type Hasher interface {
	FileHash(path, hashType string) (*hash.FileInfo, error)
	DirectoryHash(ctx context.Context, path, hashType string) error
}

type Service struct {
	Hasher
}

func NewService(repo repository.Repository) *Service {
	return &Service{Hasher: NewHasherService(repo)}
}
