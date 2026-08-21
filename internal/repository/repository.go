package repository

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"student/internal/sqlc"
)

type Repository struct {
	DB *sqlc.Queries
}

func NewRepository(db *sqlc.Queries) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) GetDB() *sqlc.Queries {
	return r.DB
}

func (r *Repository) Close() error {
	return r.DB.Close()
}

func NewDB(dns string) (*sqlc.Queries, error) {
	dir, _ := os.Getwd()
	file := filepath.Join(dir, dns)
	fmt.Printf("Connecting to database at: %s\n", file)
	conn, err := sql.Open("sqlite", file)
	if err != nil {
		return nil, err
	}
	if err = conn.Ping(); err != nil {
		return nil, err
	}
	db := sqlc.New(conn)

	return db, nil
}
