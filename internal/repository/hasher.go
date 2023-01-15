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

func (r *HasherRepository) SaveHash(input model.FileInfo, containerInfo *model.ContainerInfo) error {

	query := fmt.Sprintf(`INSERT INTO files (pod_name, image_name, image_version, 
		file_name, file_path, hash_value, hash_type) VALUES
    	($1, $2, $3, $4, $5, $6, $7) ON CONFLICT ON CONSTRAINT 
		files_unique DO UPDATE SET hash_value=excluded.hash_value`)

	_, err := r.db.Exec(query, containerInfo.PodName, containerInfo.ImageName, containerInfo.ImageVersion,
		input.FileName, input.FilePath, input.HashValue, input.HashType)

	if err != nil {
		return err
	}

	return nil
}

func (r *HasherRepository) SaveDirectoryHash(input []model.FileInfo, containerInfo *model.ContainerInfo) error {
	tx, err := r.db.Begin()

	if err != nil {
		return err
	}

	query := fmt.Sprintf(`INSERT INTO files (pod_name, image_name, image_version, 
		file_name, file_path, hash_value, hash_type) VALUES
    	($1, $2, $3, $4, $5, $6, $7) ON CONFLICT ON CONSTRAINT 
		files_unique DO UPDATE SET hash_value=excluded.hash_value`)

	for _, v := range input {
		_, err := tx.Exec(query, containerInfo.PodName, containerInfo.ImageName, containerInfo.ImageVersion,
			v.FileName, v.FilePath, v.HashValue, v.HashType)

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *HasherRepository) GetFilesInfo(dirPath, hashType string, containerInfo *model.ContainerInfo) ([]model.FileInfo,
	error) {
	var result []model.FileInfo

	query := fmt.Sprintf(`SELECT file_name, file_path, hash_value, hash_type
		FROM files WHERE file_path like $1 AND hash_type = $2 AND pod_name = $3`)

	err := r.db.Select(&result, query, "%"+dirPath+"%", hashType, containerInfo.PodName)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *HasherRepository) ClearTable(containerInfo *model.ContainerInfo) error {

	query := fmt.Sprintf(`DELETE FROM files WHERE pod_name = $1`)

	_, err := r.db.Exec(query, containerInfo.PodName)

	if err != nil {
		return err
	}

	return nil
}
