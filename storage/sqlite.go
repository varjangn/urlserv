package storage

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type SqliteStore struct {
	db *sql.DB
}

func NewSqliteStorage(dbPath string) (*SqliteStore, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &SqliteStore{
		db: db,
	}, nil
}

func (s *SqliteStore) Init() error {
	err := s.CreateUsersTable()
	return err
}

func (s *SqliteStore) CreateUsersTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        email TEXT NOT NULL UNIQUE,
        password TEXT NOT NULL,
        firstname TEXT,
		lastname TEXT
    );`
	_, err := s.db.Exec(query)
	return err
}
