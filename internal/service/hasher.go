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

func (s HasherService) CompareHash(ctx context.Context, path, hashType string) ([]model.ChangedFiles, []model.DeletedFiles, error) {
	paths := make(chan string)
	hashes := make(chan hashsum.FileInfo)

	go hashsum.Sha256sum(paths, hashes, hashType)
	go hashsum.LookUpManager(path, paths)
	newHashes := hashsum.PrintResult(ctx, hashes)

	oldHashes, err := s.repo.CompareHash(path, hashType)
	if err != nil {
		return nil, nil, err
	}

	var resultsHash []model.ChangedFiles
	var resultsDeleted []model.DeletedFiles
	var checker bool

	for _, oldHash := range oldHashes {
		checker = false
		for _, newHash := range newHashes {
			if oldHash.FilePath == newHash.FilePath {
				if oldHash.HashValue != newHash.HashValue {
					resultsHash = append(resultsHash, model.ChangedFiles{
						FileName: oldHash.FileName,
						OldHash:  oldHash.HashValue,
						NewHash:  newHash.HashValue,
					})
				}
				checker = true
				break
			}
		}
		if !checker {
			resultsDeleted = append(resultsDeleted, model.DeletedFiles{
				FileName: oldHash.FileName,
				OldHash:  oldHash.HashValue,
				FilePath: oldHash.FilePath,
			})
		}
	}

	err = s.repo.DeletedItemUpdate(resultsDeleted, hashType)

	return resultsHash, resultsDeleted, err
}
