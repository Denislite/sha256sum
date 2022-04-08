package service

type Hasher interface {
	FileHash(path, hashType string) (string, error)
	DirectoryHash()
}

type Service struct {
	Hasher
}

func NewService() *Service {
	return &Service{Hasher: NewHasherService()}
}
