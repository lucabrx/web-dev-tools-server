package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	validator "github.com/wdt/internal/validators"
)

type ToolModel struct {
	DB *sql.DB
}

type Tool struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Name      string    `json:"name"`
	Category  string    `json:"category"`
	Description string  `json:"description"`
	ImageUrl  string    `json:"imageUrl"`
	Published bool      `json:"published,omitempty"`
	Website   string    `json:"website"`
	Version   int64     `json:"version,omitempty"`
}

func ValidateTools(v *validator.Validator, tool *Tool) {
	v.Check(tool.Name != "", "name", "must be provided")
	v.Check(len(tool.Name) <= 40, "name", "must not be more than 500 bytes long")
	v.Check(tool.Category != "", "category", "must be provided")
	v.Check(len(tool.Category) <= 40, "category", "must not be more than 500 bytes long")
	v.Check(len(tool.Description) <= 160, "description", "must not be more than 5000 bytes long")
}
func (m ToolModel) Insert(tool *Tool) error {
	query := `INSERT INTO tools (name, category, image_url, description, published, website)
			 VALUES ($1, $2 ,$3, $4, $5, $6)
			 RETURNING id, created_at, version
			`

	args := []interface{}{
		tool.Name,
		tool.Category,
		NewNullString(tool.ImageUrl),
		tool.Description,
		tool.Published,
		tool.Website,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&tool.ID,
		&tool.CreatedAt,
		&tool.Version,
	)

	if err != nil {
		return err
	}

	return nil
}

func (m ToolModel) Get(id int64) (*Tool, error) {
	query := `SELECT id, created_at, name, category, coalesce(image_url, ''), description, published, version
			  FROM tools
			  WHERE id = $1`

	var tool Tool
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&tool.ID,
		&tool.CreatedAt,
		&tool.Name,
		&tool.Category,
		&tool.ImageUrl,
		&tool.Description,
		&tool.Published,
		&tool.Version,
	)
	if err != nil {
		switch {
		case err == sql.ErrNoRows:
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &tool, nil
}

func (m ToolModel) Delete(id int64) error {
	query := `DELETE FROM tools WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, id)

	if err != nil {
		return err
	}

	return nil
}

func (m ToolModel) Update(tool *Tool) error {
	query := `UPDATE tools
			  SET name = $1, category = $2, image_url = $3, description = $4, published = $5, version = version + 1
			  WHERE id = $6 AND version = $7
			  RETURNING version
			  `

	args := []interface{}{
		tool.Name,
		tool.Category,
		NewNullString(tool.ImageUrl),
		tool.Description,
		tool.Published,
		tool.ID,
		tool.Version,
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

func (m ToolModel) GetAll(filters Filters, search string) ([]*Tool, Metadata, error) {
    baseQuery := `SELECT count(*) OVER(), id, created_at, name, category, coalesce(image_url, ''), description, published, website
              FROM tools`

	advQuery := fmt.Sprintf(` ORDER BY %s %s`, filters.sortColumn(), filters.sortDirection())

    if search != "" {
        baseQuery += ` WHERE name ILIKE '%' || $3 || '%' OR category ILIKE '%' || $3 || '%'`
    }


	query := baseQuery + advQuery + ` LIMIT $1 OFFSET $2`
	args := []interface{}{filters.limit(), filters.offset()}
	if search != "" {
		args = append(args, search)
	}

    fmt.Println(query)
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    rows, err := m.DB.QueryContext(ctx, query, args...)
    if err != nil {
        return nil, Metadata{}, err
    }
    defer rows.Close()

    totalRecords := 0
    var tools []*Tool

    for rows.Next() {
        var tool Tool

        err := rows.Scan(
            &totalRecords,
            &tool.ID,
            &tool.CreatedAt,
            &tool.Name,
            &tool.Category,
            &tool.ImageUrl,
            &tool.Description,
            &tool.Published,
            &tool.Website,
        )

        if err != nil {
            return nil, Metadata{}, err
        }

        tools = append(tools, &tool)
    }

    if err = rows.Err(); err != nil {
        return nil, Metadata{}, err
    }

    metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
    return tools, metadata, nil
}


func (m ToolModel) GetAllPublished(search string, filters Filters) ([]*Tool, Metadata, error) {
	baseQuery := `SELECT count(*) OVER(), id, created_at, name, category, coalesce(image_url, ''), description, website
			  FROM tools 
			  WHERE published = true`

	if search != "" {
    	baseQuery += ` AND (name ILIKE '%' || $3 || '%' OR category ILIKE '%' || $3 || '%')`
	}

	advQuery := fmt.Sprintf(` ORDER BY %s %s`, filters.sortColumn(), filters.sortDirection())

	

	query := baseQuery + advQuery + ` LIMIT $1 OFFSET $2`


	args := []interface{}{filters.limit(), filters.offset()}

	if search != "" {
		args = append(args, search)
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	var tools []*Tool

	for rows.Next() {
		var tool Tool

		err := rows.Scan(
			&totalRecords,
			&tool.ID,
			&tool.CreatedAt,
			&tool.Name,
			&tool.Category,
			&tool.ImageUrl,
			&tool.Description,
			&tool.Website,
		)

		if err != nil {
			return nil, Metadata{},  err
		}

		tools = append(tools, &tool)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return tools, metadata,  nil
}