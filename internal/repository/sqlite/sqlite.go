package repository

import (
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
)

type Repository struct {
	DB *sql.DB
}

func NewRepository(path string) (*Repository, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	
	return &Repository{DB: db}, nil
}
