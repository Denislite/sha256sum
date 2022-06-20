package service

import (
	"k8s.io/client-go/kubernetes"
	"sha256sum/internal/model"
	"sha256sum/internal/repository"
	"sync"
	"time"
)

type Hasher interface {
	FileHash(path string) (*model.FileInfo, error)
	DirectoryHash(path string) ([]model.FileInfo, error)
	CompareHash(path string) ([]model.ChangedFiles, error)
	CheckDeleted(path string) ([]model.ChangedFiles, error)
	CheckNew(path string) ([]model.ChangedFiles, error)
	LookUpManager(inputPath string, paths chan<- string)
	Hasher(wg *sync.WaitGroup, paths <-chan string, hashes chan<- model.FileInfo)
	Sha256sum(paths chan string, hashes chan model.FileInfo)
	ReturnResult(hashes <-chan model.FileInfo) []model.FileInfo
	DirectoryCheck(ticker *time.Ticker, path string)
}

type Service struct {
	Hasher
}

func NewService(repo repository.Repository, hashType string, client *kubernetes.Clientset,
	container *model.ContainerInfo) *Service {
	return &Service{Hasher: NewHasherService(repo, hashType, client, container)}
}
