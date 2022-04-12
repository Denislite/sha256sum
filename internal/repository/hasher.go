package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"sha256sum/internal/model"
)

type HasherRepository struct {
	db *sqlx.DB
}

func NewHasherRepository(db *sqlx.DB) *HasherRepository {
	return &HasherRepository{db: db}
}

func (s *HasherRepository) SaveHash(input model.Hasher) error {
	query := fmt.Sprintf(`SELECT check_hash($1, $2, $3, $4);`)

	var id string
	row := s.db.QueryRow(query, input.FileName, input.FilePath,
		fmt.Sprintf("%x", input.HashValue), input.HashType)

	err := row.Scan(&id)
	if err != nil {
		return err
	}

	return nil
}

func (s *HasherRepository) SaveDirectoryHash(input []model.Hasher) error {
	return nil
}
