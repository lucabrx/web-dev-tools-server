package data

import (
	"github.com/stretchr/testify/require"
	validator "github.com/wdt/internal/validators"
	"testing"
)

func TestValidateTools(t *testing.T) {
	tool := &Tool{
		Name:        "test",
		Category:    "test",
		Description: "test",
	}
	v := validator.New()
	ValidateTools(v, tool)
	require.Equal(t, 0, len(v.Errors))
}

func TestToolModel_Insert(t *testing.T) {
	CreateTool(t)
}

func TestToolModel_Get_Valid(t *testing.T) {
	tool := CreateTool(t)

	dbTool, err := testQueries.Tools.Get(tool.ID)
	require.NoError(t, err)
	require.NotEmpty(t, dbTool)
	require.Equal(t, tool.ID, dbTool.ID)
	require.Equal(t, tool.Name, dbTool.Name)
	require.Equal(t, tool.Category, dbTool.Category)
	require.Equal(t, tool.Description, dbTool.Description)
}

func TestToolModel_Get_NotFound(t *testing.T) {
	_, err := testQueries.Tools.Get(0)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrRecordNotFound)
}

func TestToolModel_Delete(t *testing.T) {
	tool := CreateTool(t)

	err := testQueries.Tools.Delete(tool.ID)
	require.NoError(t, err)

	_, err = testQueries.Tools.Get(tool.ID)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrRecordNotFound)
}

func TestToolModel_Update(t *testing.T) {
	tool := CreateTool(t)

	tool.Name = "updated"
	tool.Category = "updated"
	tool.Description = "updated"

	err := testQueries.Tools.Update(&tool)
	require.NoError(t, err)

	dbTool, err := testQueries.Tools.Get(tool.ID)
	require.NoError(t, err)
	require.NotEmpty(t, dbTool)
	require.Equal(t, tool.ID, dbTool.ID)
	require.Equal(t, tool.Name, dbTool.Name)
	require.Equal(t, tool.Category, dbTool.Category)
	require.Equal(t, tool.Description, dbTool.Description)
}

func TestToolModel_GetAll(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateTool(t)
	}

	s := ""
	f := Filters{
		Page:         1,
		PageSize:     10,
		SortSafelist: []string{"id"},
		Sort:         "id",
	}

	tools, _, err := testQueries.Tools.GetAll(f, s)
	require.NoError(t, err)
	require.Len(t, tools, 10)
	require.NotEmpty(t, tools)
}

func TestToolModel_GetAllPublished(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateTool(t)
	}

	s := ""
	f := Filters{
		Page:         1,
		PageSize:     10,
		SortSafelist: []string{"id"},
		Sort:         "id",
	}

	tools, _, err := testQueries.Tools.GetAllPublished(s, f)
	require.NoError(t, err)
	require.Len(t, tools, 10)
	require.NotEmpty(t, tools)
}
