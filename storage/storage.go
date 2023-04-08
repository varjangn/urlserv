package storage

import (
	"github.com/varjangn/urlserv/types"
)

type Storage interface {
	CreateUsersTable() error
	CreateUser(u *types.User) error
	GetUserId(email string) (int64, error)
	GetUser(email string) (*types.User, error)
}
