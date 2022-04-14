package repository

import (
	"github.com/jmoiron/sqlx"
	"sha256sum/internal/model"
	"sha256sum/pkg/hashsum"
)

type Hasher interface {
	SaveHash(input hashsum.FileInfo) error
	SaveDirectoryHash(input []hashsum.FileInfo) error
	CompareHash(input []hashsum.FileInfo, dirPath string) ([]model.ChangedFiles, error)
}

type Repository struct {
	Hasher
}

func NewRepository(db *sqlx.DB) Repository {
	return Repository{Hasher: NewHasherRepository(db)}
}
