package service

import (
	"context"
	"sha256sum/internal/model"
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

	return value, err
}

func (s HasherService) DirectoryHash(ctx context.Context, path, hashType string) ([]hashsum.FileInfo, error) {
	paths := make(chan string)
	hashes := make(chan hashsum.FileInfo)

	go hashsum.Sha256sum(paths, hashes, hashType)
	go hashsum.LookUpManager(path, paths)
	result := hashsum.PrintResult(ctx, hashes)

	err := s.repo.SaveDirectoryHash(result)

	return result, err
}

func (s HasherService) CompareHash(ctx context.Context, path, hashType string) ([]model.ChangedFiles, error) {
	paths := make(chan string)
	hashes := make(chan hashsum.FileInfo)

	go hashsum.Sha256sum(paths, hashes, hashType)
	go hashsum.LookUpManager(path, paths)
	newHashes := hashsum.PrintResult(ctx, hashes)

	oldHashes, err := s.repo.GetFilesInfo(path, hashType)
	if err != nil {
		return nil, err
	}

	var resultsHash []model.ChangedFiles

	for _, oldHash := range oldHashes {
		for _, newHash := range newHashes {
			if oldHash.FilePath == newHash.FilePath && oldHash.HashValue != newHash.HashValue {
				resultsHash = append(resultsHash, model.ChangedFiles{
					FileName: oldHash.FileName,
					OldHash:  oldHash.HashValue,
					NewHash:  newHash.HashValue,
				})
			}
		}
	}

	return resultsHash, err
}

func (s HasherService) CheckDeleted(ctx context.Context, path, hashType string) ([]model.DeletedFiles, error) {
	paths := make(chan string)
	hashes := make(chan hashsum.FileInfo)

	go hashsum.Sha256sum(paths, hashes, hashType)
	go hashsum.LookUpManager(path, paths)
	newHashes := hashsum.PrintResult(ctx, hashes)

	oldHashes, err := s.repo.GetFilesInfo(path, hashType)
	if err != nil {
		return nil, err
	}

	deletedFiles := make(map[string]struct{}, len(newHashes))
	for _, value := range newHashes {
		deletedFiles[value.FilePath] = struct{}{}
	}

	var result []model.DeletedFiles

	for _, value := range oldHashes {
		if _, ok := deletedFiles[value.FilePath]; !ok {
			result = append(result, model.DeletedFiles{
				FilePath: value.FilePath,
				OldHash:  value.HashValue,
			})
		}
	}

	err = s.repo.DeletedItemUpdate(result)

	return result, err
}
