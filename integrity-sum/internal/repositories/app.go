package repositories

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/integrity-sum/internal/core/ports"
	"github.com/sirupsen/logrus"
)

type AppRepository struct {
	ports.IHashRepository
	db     *sql.DB
	logger *logrus.Logger
}

func NewAppRepository(logger *logrus.Logger, db *sql.DB) *AppRepository {
	return &AppRepository{
		IHashRepository: NewHashRepository(logger, db),
		logger:          logger,
		db:              db,
	}
}

// IsExistDeploymentNameInDB checks if the base is empty
func (ar AppRepository) IsExistDeploymentNameInDB(deploymentName string) (bool, error) {
	var count int

	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE name_deployment=$1 LIMIT 1;", os.Getenv("TABLE_NAME"))
	row := ar.db.QueryRow(query, deploymentName)
	err := row.Scan(&count)
	if err != nil {
		ar.logger.Error("err while scan row in database ", err)
		return false, err
	}

	if count < 1 {
		return true, nil
	}
	return false, nil
}
