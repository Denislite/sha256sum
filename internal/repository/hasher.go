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
		files_unique DO UPDATE SET hash_value=excluded.hash_value`)

	_, err = tx.Exec(query, input.FileName, input.FilePath, fmt.Sprintf("%x", input.HashValue), input.HashType)

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *HasherRepository) SaveDirectoryHash(input []hashsum.FileInfo) error {
	tx, err := r.db.Begin()

	if err != nil {
		return err
	}

	query := fmt.Sprintf(`INSERT INTO files (file_name, file_path, hash_value, hash_type) VALUES
    	($1, $2, $3, $4) ON CONFLICT ON CONSTRAINT 
		files_unique DO UPDATE SET hash_value=excluded.hash_value`)

	for _, v := range input {
		_, err := tx.Exec(query, v.FileName, v.FilePath, fmt.Sprintf("%x", v.HashValue), v.HashType)

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
