package repository

type Hasher interface {
	SaveHash()
}

type Repository struct {
	Hasher
}

func NewRepository() *Repository {
	return &Repository{Hasher: NewHasherRepository()}
}
