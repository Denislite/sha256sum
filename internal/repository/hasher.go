package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"sha256sum/pkg/hashsum"
)

type HasherRepository struct {
	db *sqlx.DB
}

func NewHasherRepository(db *sqlx.DB) *HasherRepository {
	return &HasherRepository{db: db}
}

func (r *HasherRepository) SaveHash(input hashsum.FileInfo) error {
	tx, err := r.db.Begin()

	if err != nil {
		return err
	}

	query := fmt.Sprintf(`INSERT INTO files (file_name, file_path, hash_value, hash_type) VALUES
    	($1, $2, $3, $4) ON CONFLICT ON CONSTRAINT 
		files_unique DO UPDATE SET hash_value=excluded.hash_value RETURNING id;`)

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

func (r *HasherRepository) SaveDirectoryHash(input []hashsum.FileInfo) error {
	for _, hash := range input {
		r.SaveHash(hash)
	}
	return nil
}
