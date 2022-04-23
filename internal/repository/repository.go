package repository

import (
	"github.com/jmoiron/sqlx"
	"sha256sum/internal/model"
)

type Hasher interface {
	SaveHash(input model.FileInfo) error
	SaveDirectoryHash(input []model.FileInfo) error
	GetFilesInfo(dirPath, hashType string) ([]model.FileInfo, error)
	DeletedItemUpdate(input []model.DeletedFiles) error
}

type Repository struct {
	Hasher
}

func NewRepository(db *sqlx.DB) Repository {
	return Repository{Hasher: NewHasherRepository(db)}
}
