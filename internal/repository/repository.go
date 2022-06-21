package repository

import (
	"github.com/jmoiron/sqlx"
	"sha256sum/internal/model"
)

type Hasher interface {
	SaveHash(input model.FileInfo, containerInfo *model.ContainerInfo) error
	SaveDirectoryHash(input []model.FileInfo, containerInfo *model.ContainerInfo) error
	GetFilesInfo(dirPath, hashType string, containerInfo *model.ContainerInfo) ([]model.FileInfo, error)
	ClearTable(containerInfo *model.ContainerInfo) error
}

type Repository struct {
	Hasher
}

func NewRepository(db *sqlx.DB) Repository {
	return Repository{Hasher: NewHasherRepository(db)}
}
