package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"sha256sum/internal/model"
	"sha256sum/pkg/hashsum"
)

type HasherRepository struct {
	db *sqlx.DB
}

func NewHasherRepository(db *sqlx.DB) *HasherRepository {
	return &HasherRepository{db: db}
}

func (r *HasherRepository) SaveHash(input hashsum.FileInfo) error {

	query := fmt.Sprintf(`INSERT INTO files (file_name, file_path, hash_value, hash_type) VALUES
    	($1, $2, $3, $4) ON CONFLICT ON CONSTRAINT 
		files_unique DO UPDATE SET hash_value=excluded.hash_value`)

	_, err := r.db.Exec(query, input.FileName, input.FilePath, input.HashValue, input.HashType)

	if err != nil {
		return err
	}

	return nil
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
		_, err := tx.Exec(query, v.FileName, v.FilePath, v.HashValue, v.HashType)

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *HasherRepository) CompareHash(dirPath, hashType string) ([]hashsum.FileInfo, error) {
	var result []hashsum.FileInfo

	query := fmt.Sprintf(`SELECT file_name, file_path, hash_value, hash_type, deleted 
		FROM files WHERE file_path like $1 AND hash_type = $2`)

	err := r.db.Select(&result, query, "%"+dirPath+"%", hashType)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *HasherRepository) DeletedItemUpdate(input []model.DeletedFiles, hashType string) error {
	tx, err := r.db.Begin()

	if err != nil {
		return err
	}

	query := fmt.Sprintf(`UPDATE files SET deleted = true WHERE file_path=$1 AND hash_type=$2`)

	for _, v := range input {
		_, err := tx.Exec(query, v.FilePath, hashType)

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}