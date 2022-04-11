package repository

import (
	"github.com/jmoiron/sqlx"
	"sha256sum/internal/model"
)

type Hasher interface {
	SaveHash(input model.Hasher) error
	SaveDirectoryHash(input []model.Hasher) error
}

type Repository struct {
	Hasher
}

func NewRepository(db *sqlx.DB) Repository {
	return Repository{Hasher: NewHasherRepository(db)}
}
