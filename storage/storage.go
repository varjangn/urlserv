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
}
