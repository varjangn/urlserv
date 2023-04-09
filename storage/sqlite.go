package storage

import (
	"database/sql"
	"fmt"
	"time"

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
	if err != nil {
		return err
	}
	err = s.CreateUrlsTable()
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

func (s *SqliteStore) CreateUrlsTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS urls (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,
        short_id TEXT NOT NULL UNIQUE,
		long TEXT NOT NULL,
		created_at INTEGER NOT NULL
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

func (s *SqliteStore) CreateURL(url *types.URL) error {
	createdAt := time.Now().UTC().UnixMilli()
	qry := "INSERT INTO urls(user_id, short_id, long, created_at) values(?,?,?,?)"
	res, err := s.db.Exec(qry, url.UserId, url.ShortId, url.Long, createdAt)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil
	}
	url.Id = id
	url.CreatedAt = createdAt
	return nil
}

func (s *SqliteStore) GetLongURL(shortId string) (string, error) {
	qry := fmt.Sprintf("SELECT long FROM urls WHERE short_id='%s'", shortId)
	row := s.db.QueryRow(qry)
	var longURL string

	switch err := row.Scan(&longURL); err {
	case sql.ErrNoRows:
		return longURL, nil
	case nil:
		return longURL, nil
	default:
		return longURL, err
	}
}

func (s *SqliteStore) GetURLbyLongURL(longURL string) (*types.URL, error) {
	qry := fmt.Sprintf("SELECT * FROM urls WHERE long='%s'", longURL)
	row := s.db.QueryRow(qry)
	url := new(types.URL)

	switch err := row.Scan(&url.Id, &url.UserId, &url.ShortId, &url.Long, &url.CreatedAt); err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		return url, nil
	default:
		return nil, err
	}
}

func (s *SqliteStore) GetURLs(u *types.User) ([]*types.URL, error) {
	qry := fmt.Sprintf("SELECT * FROM urls WHERE user_id=%d", u.Id)
	rows, err := s.db.Query(qry)
	if err != nil {
		return nil, err
	}
	urls := []*types.URL{}
	for rows.Next() {
		url := new(types.URL)
		err = rows.Scan(
			&url.Id, &url.UserId, &url.ShortId, &url.Long,
			&url.CreatedAt)
		if err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}
	return urls, nil
}
