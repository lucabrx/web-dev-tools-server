package data

import (
	"github.com/stretchr/testify/require"
	"github.com/wdt/internal/random"
	validator "github.com/wdt/internal/validators"
	"testing"
)

func TestUserModel_Insert_Valid(t *testing.T) {
	CreateRandomUser(t)
}

func TestUserModel_Insert_UsedEmail(t *testing.T) {
	email := random.RandString(10) + "@gmail.com"
	user1 := &User{
		Email: email,
		Name:  random.RandString(10),
	}
	user2 := &User{
		Email: email,
		Name:  random.RandString(10),
	}

	err := testQueries.Users.Insert(user1)
	require.NoError(t, err)
	require.Regexp(t, validator.EmailRX, user1.Email)
	require.NotZero(t, user1.ID)
	require.NotZero(t, user1.CreatedAt)
	require.NotZero(t, user1.Version)

	err = testQueries.Users.Insert(user2)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrDuplicateEmail)
}

func TestUserModel_Get_Valid_ID(t *testing.T) {
	user := CreateRandomUser(t)

	dbUser, err := testQueries.Users.Get(user.ID, "")
	require.NoError(t, err)

	require.Equal(t, user.ID, dbUser.ID)
	require.Equal(t, user.Email, dbUser.Email)
	require.Equal(t, user.Name, dbUser.Name)
}
func TestUserModel_Get_Valid_Email(t *testing.T) {
	user := CreateRandomUser(t)

	dbUser, err := testQueries.Users.Get(0, user.Email)
	require.NoError(t, err)

	require.Equal(t, user.ID, dbUser.ID)
	require.Equal(t, user.Email, dbUser.Email)
	require.Equal(t, user.Name, dbUser.Name)
}
func TestUserModel_Get_Invalid(t *testing.T) {
	_, err := testQueries.Users.Get(0, "")
	require.Error(t, err)
	require.ErrorIs(t, err, ErrRecordNotFound)
}

func TestUserModel_GetForToken(t *testing.T) {
	user := CreateRandomUser(t)
	token := CreateTokenForUser(t, user)

	dbUser, err := testQueries.Users.GetForToken(ScopeAuthentication, token.Plaintext)
	require.NoError(t, err)

	require.Equal(t, user.ID, dbUser.ID)
	require.Equal(t, user.Email, dbUser.Email)
	require.Equal(t, user.Name, dbUser.Name)
}

func TestUserModel_GetForToken_NotFound(t *testing.T) {
	_, err := testQueries.Users.GetForToken(ScopeAuthentication, "")
	require.Error(t, err)
	require.ErrorIs(t, err, ErrRecordNotFound)
}
