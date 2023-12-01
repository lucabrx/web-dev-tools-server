package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	validator "github.com/wdt/internal/validators"
)

type CategoryModel struct {
	DB *sql.DB
}

type Category struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	Name      string    `json:"name"`
	Published bool      `json:"published"`
	Version   int64     `json:"version,omitempty"`
}

func ValidateCategories(v *validator.Validator, category *Category) {
	v.Check(category.Name != "", "name", "must be provided")
	v.Check(len(category.Name) <= 500, "name", "must not be more than 500 bytes long")
}

func (m CategoryModel) Insert(category *Category) error {
	query := `INSERT INTO categories (name)
			 VALUES ($1)
			 RETURNING id, created_at, version
			`

	args := []interface{}{
		category.Name,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&category.ID,
		&category.CreatedAt,
		&category.Version,
	)

	if err != nil {
		return err
	}

	return nil
}

func (m CategoryModel) Get(id int64) (*Category, error) {
	query := `SELECT id, created_at, name, published, version
			  FROM categories
			  WHERE id = $1`

	var category Category

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&category.ID,
		&category.CreatedAt,
		&category.Name,
		&category.Published,
		&category.Version,
	)

	if err != nil {

		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &category, nil
}

func (m CategoryModel) GetAll(filters Filters) ([]*Category, Metadata, error) {
	query := fmt.Sprintf(`SELECT count(*) OVER(), id, name, published, created_at
			  FROM categories
			  ORDER BY %s %s
			  LIMIT $1 OFFSET $2`,
		filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, filters.limit(), filters.offset())
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	var categories []*Category

	for rows.Next() {
		var category Category
		err := rows.Scan(
			&totalRecords,
			&category.ID,
			&category.Name,
			&category.Published,
			&category.CreatedAt,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		categories = append(categories, &category)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	return categories, metadata, nil
}

func (m CategoryModel) Update(category *Category) error {
	query := `UPDATE categories SET name = $1, published = $2, version = version + 1 WHERE id = $3`

	args := []interface{}{
		category.Name,
		category.Published,
		category.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, args...)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (m CategoryModel) Delete(id int64) error {
	query := `DELETE FROM categories WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, id)

	if err != nil {
		return err
	}

	return nil
}

func (m CategoryModel) GetAllPublished() ([]*Category, error) {
	query := `SELECT id, name
			  FROM categories
			  WHERE published = true
			  ORDER BY name
			 `

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := []*Category{}

	for rows.Next() {
		var category Category
		err := rows.Scan(
			&category.ID,
			&category.Name,
		)
		if err != nil {
			return nil, err
		}

		categories = append(categories, &category)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}
