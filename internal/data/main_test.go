package data

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"github.com/stretchr/testify/require"
	"github.com/wdt/internal/random"
	validator "github.com/wdt/internal/validators"
	"time"

	"testing"

	_ "github.com/lib/pq"
)

var testQueries Models

func TestMain(m *testing.M) {
	conn, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		panic(err)
	}

	testQueries = NewModels(conn)

	m.Run()
}

func CreateRandomUser(t *testing.T) User {
	user := &User{
		Email: random.RandString(10) + "@gmail.com",
		Name:  random.RandString(10),
	}

	err := testQueries.Users.Insert(user)
	require.NoError(t, err)
	require.NotEmptyf(t, user.Email, "user email should not be empty")

	require.Regexp(t, validator.EmailRX, user.Email)

	require.NotZero(t, user.ID)
	require.NotZero(t, user.CreatedAt)
	require.NotZero(t, user.Version)

	return *user
}

func CreateTokenForUser(t *testing.T, user User) Token {
	token := &Token{
		UserID: user.ID,
		Expiry: time.Now().Add(24 * time.Hour),
		Scope:  ScopeAuthentication,
	}
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	require.NoError(t, err)
	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]

	err = testQueries.Tokens.Insert(token)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	return *token
}

func CreateCategory(t *testing.T) Category {
	category := Category{
		Name:      random.RandString(10),
		Published: true,
	}

	err := testQueries.Categories.Insert(&category)
	require.NoError(t, err)
	require.NotEmpty(t, category)

	return category
}

func CreateTool(t *testing.T) Tool {
	tool := &Tool{
		Name:        random.RandString(10),
		Category:    random.RandString(5),
		Description: random.RandString(20),
	}

	err := testQueries.Tools.Insert(tool)
	require.NoError(t, err)
	require.NotZero(t, tool.ID)
	require.NotZero(t, tool.CreatedAt)
	require.NotZero(t, tool.Version)

	return *tool
}
