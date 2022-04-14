package service

import (
	"context"
	"sha256sum/internal/model"
	"sha256sum/internal/repository"
	"sha256sum/pkg/hashsum"
)

type Hasher interface {
	FileHash(path, hashType string) (*hashsum.FileInfo, error)
	DirectoryHash(ctx context.Context, path, hashType string, check bool) ([]model.ChangedFiles, error)
}

type Service struct {
	Hasher
}

func NewService(repo repository.Repository) *Service {
	return &Service{Hasher: NewHasherService(repo)}
}
