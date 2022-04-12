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

func (r *HasherRepository) SaveHash(input model.Hasher) error {
	tx, err := r.db.Begin()

	if err != nil {
		return err
	}

	query := fmt.Sprintf(`SELECT check_hash($1, $2, $3, $4);`)

	var id string
	row := tx.QueryRow(query, input.FileName, input.FilePath,
		fmt.Sprintf("%x", input.HashValue), input.HashType)

	err = row.Scan(&id)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *HasherRepository) SaveDirectoryHash(input []model.Hasher) error {
	//just for test
	//for _, hash := range input {
	//	r.SaveHash(hash)
	//}
	return nil
}
