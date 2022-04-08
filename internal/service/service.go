package service

type Hasher interface {
	FileHash()
	DirectoryHash()
}

type Service struct {
	Hasher
}

func NewService() *Service {
	return &Service{Hasher: NewHasherService()}
}
