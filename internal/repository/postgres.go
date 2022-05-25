package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"os"
)

func NewPostgresDB() (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres",
		fmt.Sprintf("postgres://%s:%s@%v:%v/%s?sslmode=%s", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"), os.Getenv("DB_SSL")))

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
