package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Movies MovieModel
	Users  UserModel
	Tokens TokenModel
}

func NewModels(db *sql.DB) Models {
	return Models{Movies: MovieModel{DB: db},
		Users:  UserModel{Db: db},
		Tokens: TokenModel{Db: db},
	}
}
