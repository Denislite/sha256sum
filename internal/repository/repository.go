package repository

import (
	"github.com/jmoiron/sqlx"
	"sha256sum/pkg/hash"
)

type Hasher interface {
	SaveHash(input hash.FileInfo) error
	SaveDirectoryHash(input []hash.FileInfo) error
}

type Repository struct {
	Hasher
}

func NewRepository(db *sqlx.DB) Repository {
	return Repository{Hasher: NewHasherRepository(db)}
}
