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
	Users      UserModel
	Tokens     TokenModel
	Tools      ToolModel
	Categories CategoryModel
	Favorites  FavoriteModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Users:      UserModel{DB: db},
		Tokens:     TokenModel{DB: db},
		Tools:      ToolModel{DB: db},
		Categories: CategoryModel{DB: db},
		Favorites:  FavoriteModel{DB: db},
	}
}

func NewNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}
