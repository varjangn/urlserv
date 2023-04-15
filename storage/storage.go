package storage

import (
	"github.com/varjangn/urlserv/types"
)

type Storage interface {
	CreateUsersTable() error
	CreateUser(u *types.User) error
	GetUserId(email string) (int64, error)
	GetUser(email string) (*types.User, error)
	CreateUrlsTable() error
	CreateURL(url *types.URL) error
	GetLongURL(shortId string) (string, error)
	GetURLbyLongURL(longURL string) (*types.URL, error)
	GetURLs(u *types.User) ([]*types.URL, error)
	GetURL(id int64, userId int64) (*types.URL, error)
	DeleteURL(id int64, userId int64) (bool, error)
	UpdateLongURL(id int64, userId int64, LongURL string) (bool, error)
}
