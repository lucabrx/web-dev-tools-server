package data

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTokenModel_Insert(t *testing.T) {
	user := CreateRandomUser(t)
	CreateTokenForUser(t, user)
}

func TestTokenModel_DeleteAllForUser(t *testing.T) {
	user := CreateRandomUser(t)
	CreateTokenForUser(t, user)

	err := testQueries.Tokens.DeleteAllForUser(ScopeAuthentication, user.ID)
	require.NoError(t, err)

}
