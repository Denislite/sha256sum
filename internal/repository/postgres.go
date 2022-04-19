package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"sha256sum/internal/model"
)

func NewPostgresDB(cfg *model.PostgresConfig) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres",
		fmt.Sprintf("postgres://%s:%s@%v:%v/%s?sslmode=%s",
			cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode))

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
