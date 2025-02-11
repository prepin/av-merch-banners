package entities

import "golang.org/x/crypto/bcrypt"

type User struct {
	ID             int
	Username       string
	HashedPassword string
	Role           string
}

func NewUserWithPassword(username string, password string) (User, error) {
	u := User{Username: username}
	err := u.SetPassword(password)
	return u, err
}

func (u *User) SetPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.HashedPassword = string(bytes)
	return nil
}

func (u *User) CompareWithPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(password))
	return err == nil
}
