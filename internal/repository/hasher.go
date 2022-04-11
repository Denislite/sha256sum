package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
)

type HasherRepository struct {
	db *sqlx.DB
}

func NewHasherRepository(db *sqlx.DB) *HasherRepository {
	return &HasherRepository{db: db}
}

func (s *HasherRepository) SaveHash(name, hash string) error {
	query := fmt.Sprintf("INSERT INTO files (file_name, hash_value) VALUES ($1, $2) RETURNING id")

	var id string
	row := s.db.QueryRow(query, name, hash)

	err := row.Scan(&id)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
