package data

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCategoryModel_Insert(t *testing.T) {
	CreateCategory(t)
}

func TestCategoryModel_Get_Valid(t *testing.T) {
	category := CreateCategory(t)

	dbCategory, err := testQueries.Categories.Get(category.ID)
	require.NoError(t, err)
	require.NotEmpty(t, dbCategory)
	require.Equal(t, category.ID, dbCategory.ID)
	require.Equal(t, category.Name, dbCategory.Name)
}

func TestCategoryModel_Get_NotFound(t *testing.T) {
	_, err := testQueries.Categories.Get(0)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrRecordNotFound)
}

func TestCategoryModel_GetAll(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateCategory(t)
	}

	f := Filters{
		SortSafelist: []string{"id"},
		PageSize:     20,
		Page:         1,
		Sort:         "id",
	}

	categories, _, err := testQueries.Categories.GetAll(f)
	require.NoError(t, err)
	require.NotEmpty(t, categories)

	for _, category := range categories {
		require.NotEmpty(t, category)
	}
}

func TestCategoryModel_Update(t *testing.T) {
	category := CreateCategory(t)

	category.Name = "Updated Name"

	err := testQueries.Categories.Update(&category)
	require.NoError(t, err)

	dbCategory, err := testQueries.Categories.Get(category.ID)
	require.NoError(t, err)
	require.NotEmpty(t, dbCategory)
	require.Equal(t, category.ID, dbCategory.ID)
	require.Equal(t, category.Name, dbCategory.Name)
}

func TestCategoryModel_Delete(t *testing.T) {
	category := CreateCategory(t)

	err := testQueries.Categories.Delete(category.ID)
	require.NoError(t, err)

	dbCategory, err := testQueries.Categories.Get(category.ID)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrRecordNotFound)
	require.Empty(t, dbCategory)
}

func TestCategoryModel_GetAllPublished(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateCategory(t)
	}

	categories, err := testQueries.Categories.GetAllPublished()
	require.NoError(t, err)
	require.NotEmpty(t, categories)

	for _, category := range categories {
		require.NotEmpty(t, category)
	}
}
