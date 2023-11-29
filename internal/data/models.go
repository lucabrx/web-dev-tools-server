package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound   = errors.New("record not found")
	ErrEditConflict     = errors.New("edit conflict")
	ErrFieldRequired    = errors.New("field required")
	ErrInvalidEnumValue = errors.New("invalid enum value")
)

type Models struct {
	Users       UserModel
	Tokens      TokenModel
	Tools 		ToolModel
	Categories  CategoryModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Users: UserModel{DB: db},
		Tokens: TokenModel{DB: db},
		Tools: ToolModel{DB: db},
		Categories: CategoryModel{DB: db},
	}
}

func sqlArray(s []string) string {
	var result string
	for i, v := range s {
		if i == len(s)-1 {
			result += v
		} else {
			result += v + ","
		}
	}
	return result
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

func NewNullPgArray(s []string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: "{" + sqlArray(s) + "}",
		Valid:  true,
	}

}

func NewNullNil(s []byte) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: string(s),
		Valid:  true,
	}
}