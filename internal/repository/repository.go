package repository

type Hasher interface {
	SaveHash(name, hash string) error
}

type Repository struct {
	Hasher
}

func NewRepository() *Repository {
	return &Repository{Hasher: NewHasherRepository()}
}
