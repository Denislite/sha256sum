package service

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sha256sum/internal/model"
	"sha256sum/internal/repository"
	"sha256sum/internal/utils"
	"sha256sum/pkg/hashsum"
	"sync"
	"time"
)

type HasherService struct {
	repo          repository.Repository
	hashSum       hashsum.HashSum
	hashType      string
	clientSet     *kubernetes.Clientset
	containerInfo *model.ContainerInfo
}

func NewHasherService(repo repository.Repository, hashType string, client *kubernetes.Clientset,
	container *model.ContainerInfo) *HasherService {
	h, t, err := hashsum.New(hashType)

	if err != nil {
		return nil
	}

	return &HasherService{
		repo:          repo,
		hashSum:       h,
		hashType:      t,
		clientSet:     client,
		containerInfo: container,
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

func (s HasherService) DirectoryHash(path string) ([]model.FileInfo, error) {
	dbHashes, err := s.repo.GetFilesInfo(path, s.hashType, s.containerInfo)

	if err != nil {
		return nil, err
	}

	if dbHashes != nil {
		return dbHashes, nil
	}

	paths := make(chan string)
	hashes := make(chan model.FileInfo)

	go s.Sha256sum(paths, hashes)
	go s.LookUpManager(path, paths)
	result := s.ReturnResult(hashes)

	return result, s.repo.SaveDirectoryHash(result, s.containerInfo)
}

func (s HasherService) CompareHash(path string) ([]model.ChangedFiles, error) {
	paths := make(chan string)
	hashes := make(chan model.FileInfo)

	go s.Sha256sum(paths, hashes)
	go s.LookUpManager(path, paths)
	newHashes := s.ReturnResult(hashes)

	oldHashes, err := s.repo.GetFilesInfo(path, s.hashType, s.containerInfo)
	if err != nil {
		return nil, err
	}

	var resultsHash []model.ChangedFiles

	resultsHash = append(resultsHash, s.CheckFiles(newHashes, oldHashes)...)

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

	resultsHash = append(resultsHash, s.CheckFiles(oldHashes, newHashes)...)

	return resultsHash, err
}

func (s HasherService) CheckFiles(newHashes, oldHashes []model.FileInfo) []model.ChangedFiles {
	Files := make(map[string]struct{}, len(newHashes))
	for _, value := range oldHashes {
		Files[value.FilePath] = struct{}{}
	}

	var result []model.ChangedFiles

	for _, value := range newHashes {
		if _, ok := Files[value.FilePath]; !ok {
			result = append(result, model.ChangedFiles{
				FileName: value.FileName,
				OldHash:  value.HashValue,
			})
		}
	}

	return result
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

func (s HasherService) ReturnResult(hashes <-chan model.FileInfo) []model.FileInfo {
	var result []model.FileInfo
	for {
		select {
		case hash, ok := <-hashes:
			if !ok {
				return result
			}
			result = append(result, hash)
		}
	}
}

func (s HasherService) DirectoryCheck(ticker *time.Ticker, path string) {

	patchData := fmt.Sprintf(`{"spec":{"template":{"metadata":{"annotations":{"kubectl.kubernetes.io/restartedAt":"%s"}}}}}`, time.Now().Format(time.RFC3339))

	for {
		select {
		case <-ticker.C:
			result, err := s.CompareHash(path)
			if err != nil {
				log.Fatalln(err)
			}
			if result != nil {
				log.Println("=========================")
				log.Println("### ❌  Files was changed:")

				for _, hash := range result {
					log.Printf("### %s %s %s \n",
						hash.FileName, hash.OldHash, hash.NewHash)
				}

				_, err = s.clientSet.AppsV1().Deployments(os.Getenv("NAMESPACE")).Patch(context.Background(),
					os.Getenv("DEPLOYMENT_NAME"), types.StrategicMergePatchType, []byte(patchData),
					metav1.PatchOptions{FieldManager: "kubectl-rollout"})

				if err != nil {
					log.Fatalf("### 👎 Warning: Failed to patch %s, restart failed: %v",
						"deployment", err)
				} else {
					err = s.repo.ClearTable(s.containerInfo)
					if err != nil {
						log.Println(err)
					}
					log.Fatalf("### ✅ Target %s, named %s was restarted!",
						"deployment", os.Getenv("DEPLOYMENT_NAME"))
				}
			}
			log.Println("### ✅  Directory was checked, all right")
		}
	}
}
