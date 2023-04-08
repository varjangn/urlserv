package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/varjangn/urlserv/types"
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
		verified BOOLEAN NOT NULL CHECK (verified IN (0, 1)),
        first_name TEXT,
		last_name TEXT
    );`
	_, err := s.db.Exec(query)
	return err
}

func (s *SqliteStore) CreateUser(u *types.User) error {
	qry := "INSERT INTO users(email, password, verified, first_name, last_name) values(?,?,?,?,?)"
	res, err := s.db.Exec(qry, u.Email, u.Password, u.Verified, u.FirstName, u.LastName)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	u.Id = id
	return nil
}

func (s *SqliteStore) GetUserId(email string) (int64, error) {
	qry := fmt.Sprintf("SELECT id from users WHERE email='%s'", email)
	row := s.db.QueryRow(qry)
	var userId int64 = -1

	switch err := row.Scan(&userId); err {
	case sql.ErrNoRows:
		return userId, nil
	case nil:
		return userId, nil
	default:
		return userId, err
	}
}

func (s *SqliteStore) GetUser(email string) (*types.User, error) {
	qry := fmt.Sprintf("SELECT * from users WHERE email='%s'", email)
	row := s.db.QueryRow(qry)
	var user types.User

	switch err := row.Scan(&user.Id, &user.Email, &user.Password, &user.Verified,
		&user.FirstName, &user.LastName); err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		return &user, nil
	default:
		return nil, err
	}
}
