package service

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sha256sum/internal/model"
	"sha256sum/internal/repository"
	"sha256sum/internal/utils"
	"sha256sum/pkg/hashsum"
	"sync"
)

type HasherService struct {
	repo     repository.Repository
	hashSum  hashsum.HashSum
	hashType string
}

func NewHasherService(repo repository.Repository, hashType string) *HasherService {
	h, t, err := hashsum.New(hashType)

	if err != nil {
		return nil
	}

	return &HasherService{
		repo:     repo,
		hashSum:  h,
		hashType: t,
	}
}

func (s HasherService) FileHash(path string) (*model.FileInfo, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, utils.ErrorWrongFile
	}

	defer file.Close()

	result, err := s.hashSum.Hash(file)

	if err != nil {
		return nil, utils.ErrorHash
	}

	data := model.FileInfo{
		FileName:  filepath.Base(path),
		FilePath:  path,
		HashType:  s.hashType,
		HashValue: result,
	}

	return &data, nil
}

func (s HasherService) DirectoryHash(ctx context.Context, path string) ([]model.FileInfo, error) {
	paths := make(chan string)
	hashes := make(chan model.FileInfo)

	go s.Sha256sum(paths, hashes)
	go s.LookUpManager(path, paths)
	result := s.ReturnResult(ctx, hashes)

	err := s.repo.SaveDirectoryHash(result)

	return result, err
}

func (s HasherService) CompareHash(ctx context.Context, path string) ([]model.ChangedFiles, error) {
	paths := make(chan string)
	hashes := make(chan model.FileInfo)

	go s.Sha256sum(paths, hashes)
	go s.LookUpManager(path, paths)
	newHashes := s.ReturnResult(ctx, hashes)

	oldHashes, err := s.repo.GetFilesInfo(path, s.hashType)
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

func (s HasherService) CheckDeleted(ctx context.Context, path string) ([]model.DeletedFiles, error) {
	paths := make(chan string)
	hashes := make(chan model.FileInfo)

	go s.Sha256sum(paths, hashes)
	go s.LookUpManager(path, paths)
	newHashes := s.ReturnResult(ctx, hashes)

	oldHashes, err := s.repo.GetFilesInfo(path, s.hashType)
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

func (s HasherService) LookUpManager(inputPath string, paths chan<- string) {
	err := filepath.Walk(inputPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return utils.ErrorDirectoryRead
		}
		if !info.IsDir() {
			paths <- path
		}

		return nil
	})
	close(paths)

	if err != nil {
		log.Println(err)
		return
	}
}

func (s HasherService) Hasher(wg *sync.WaitGroup, paths <-chan string, hashes chan<- model.FileInfo) {
	defer wg.Done()
	for path := range paths {
		hash, err := s.FileHash(path)
		if err != nil {
			log.Println(err)
		}
		hashes <- *hash
	}
}

func (s HasherService) Sha256sum(paths chan string, hashes chan model.FileInfo) {
	var wg sync.WaitGroup
	for worker := 1; worker <= runtime.NumCPU(); worker++ {
		wg.Add(1)
		go s.Hasher(&wg, paths, hashes)
	}
	defer close(hashes)
	wg.Wait()
}

func (s HasherService) ReturnResult(ctx context.Context, hashes <-chan model.FileInfo) []model.FileInfo {
	var result []model.FileInfo
	for {
		select {
		case hash, ok := <-hashes:
			if !ok {
				return result
			}
			result = append(result, hash)
			// modified for k8s deployment
			//case <-ctx.Done():
			//	log.Println("request canceled by context")
			//	os.Exit(1)
			//	return nil
		}
	}
	return result
}
