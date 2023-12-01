package data

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFavoriteModel_AddFavorite(t *testing.T) {
	user := CreateRandomUser(t)
	tool := CreateTool(t)

	err := testQueries.Favorites.AddFavorite(user.ID, tool.ID)
	require.NoError(t, err)
}

func TestFavoriteModel_DeleteFavorite(t *testing.T) {
	user := CreateRandomUser(t)
	tool := CreateTool(t)

	err := testQueries.Favorites.AddFavorite(user.ID, tool.ID)
	require.NoError(t, err)

	err = testQueries.Favorites.RemoveFavorite(user.ID, tool.ID)
	require.NoError(t, err)
}

func TestFavoriteModel_GetFavorites(t *testing.T) {
	user := CreateRandomUser(t)
	tool := CreateTool(t)

	err := testQueries.Favorites.AddFavorite(user.ID, tool.ID)
	require.NoError(t, err)

	favorites, err := testQueries.Favorites.GetFavorites(user.ID)
	require.NoError(t, err)
	require.NotEmpty(t, favorites)

	for _, favorite := range favorites {
		require.NotEmpty(t, favorite)
		require.Equal(t, user.ID, favorite.UserId)
		require.Equal(t, tool.ID, favorite.ToolId)
	}
}
