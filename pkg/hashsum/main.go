package hashsum

type Hasher interface {
	FileHash()
}

func New(hashType string) (hasher Hasher, err error) {
	//TODO implement me
	return nil, nil
}
