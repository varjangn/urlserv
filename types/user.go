package types

import "golang.org/x/crypto/bcrypt"

type User struct {
	Id        int64  `json:"id"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	Verified  bool   `json:"is_verified"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

func NewUser(email, password, firstname, lastname string) (*User, error) {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return nil, err
	}
	return &User{
		Id:        0,
		Email:     email,
		Password:  string(hashedPwd),
		Verified:  false,
		FirstName: firstname,
		LastName:  lastname,
	}, nil
}

func (u *User) UpdatePassword(newPassword string) (string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(newPassword), 10)
	if err != nil {
		return "", err
	}
	u.Password = string(hashedPass)
	return string(hashedPass), nil
}
