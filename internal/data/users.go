package data

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"time"

	validator "github.com/wdt/internal/validators"
)


var ErrDuplicateEmail = errors.New("duplicate email")
var AnonymousUser = &User{}

type UserModel struct {
	DB *sql.DB
}

type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name,omitempty"`
	Email     string    `json:"email"`
	ImageUrl  string    `json:"imageUrl,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	Version   int64     `json:"version,omitempty"`
	Role      string    `json:"role"`
}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "Email must be provided.")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address.")
}


func (m UserModel) Insert(user *User) error {
	query := `INSERT INTO users (name, email, image_url)
			 VALUES ($1, $2 ,$3)
			 RETURNING id, created_at, version
			`

	args := []interface{}{
		user.Name,
		user.Email,
		NewNullString(user.ImageUrl),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Version,
	)
	if err != nil {
		switch {
		case err.Error() == "pq: duplicate key value violates unique constraint \"users_email_key\"":
			return ErrDuplicateEmail
		default:
			return err
		}
	}

	return nil
}

func (m UserModel) Get(id int64, email string) (*User, error) {
	query := `SELECT id, created_at, COALESCE(name, ''), email, COALESCE(image_url, ''), version, role
			 FROM users
			 WHERE id = $1 OR email = $2`

	var user User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id, email).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Email,
		&user.ImageUrl,
		&user.Version,
		&user.Role,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (m UserModel) GetForToken(tokenScope, tokenPlaintext string) (*User, error) {
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))
	query := `
        SELECT u.id, u.created_at, COALESCE(u.name, ''), u.email ,COALESCE(u.image_url,''), u.version, u.role
        FROM users u
        INNER JOIN tokens t
        ON u.id = t.user_id
        WHERE t.hash = $1
        AND t.scope = $2 
        AND t.expiry > $3`

	args := []any{tokenHash[:], tokenScope, time.Now()}

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Email,
		&user.ImageUrl,
		&user.Version,
		&user.Role,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}