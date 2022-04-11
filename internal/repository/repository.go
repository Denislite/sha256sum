package repository

import "github.com/jmoiron/sqlx"

type Hasher interface {
	SaveHash(name, hash string) error
}

type Repository struct {
	Hasher
}

func NewRepository(db *sqlx.DB) Repository {
	return Repository{Hasher: NewHasherRepository(db)}
}
