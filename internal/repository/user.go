package repository

import (
	"av-merch-shop/internal/entities"
	"av-merch-shop/internal/errs"
	"av-merch-shop/pkg/database"
	"errors"
)

type PGUserRepo struct {
	db *database.Database
}

func NewPGUserRepo(db *database.Database) *PGUserRepo {
	return &PGUserRepo{
		db: db,
	}
}

func (u *PGUserRepo) GetByUsername(username string) (*entities.User, error) {
	if username == "employee" {
		user, err := entities.NewUserWithPassword(username, "password")
		return &user, err
	}
	if username == "bob" {
		return nil, errors.New("db broken")
	}
	return &entities.User{}, errs.ErrNotFound{Err: errors.New("not found")}
}

func (u *PGUserRepo) Create() {

}
