package service

import (
	"context"
	"sha256sum/internal/model"
	"sha256sum/internal/repository"
	"sync"
)

type Hasher interface {
	FileHash(path string) (*model.FileInfo, error)
	DirectoryHash(ctx context.Context, path string) ([]model.FileInfo, error)
	CompareHash(ctx context.Context, path string) ([]model.ChangedFiles, error)
	CheckDeleted(ctx context.Context, path string) ([]model.DeletedFiles, error)
	LookUpManager(inputPath string, paths chan<- string)
	Hasher(wg *sync.WaitGroup, paths <-chan string, hashes chan<- model.FileInfo)
	Sha256sum(paths chan string, hashes chan model.FileInfo)
	ReturnResult(ctx context.Context, hashes <-chan model.FileInfo) []model.FileInfo
}

type Service struct {
	Hasher
}

func NewService(repo repository.Repository, hashType string) *Service {
	return &Service{Hasher: NewHasherService(repo, hashType)}
}
