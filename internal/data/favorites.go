package data

import (
	"context"
	"database/sql"
	"time"
)

type FavoriteModel struct {
	DB *sql.DB
}

type Favorite struct {
	UserId int64 `json:"user_id"`
	ToolId int64 `json:"tool_id"`
}

func (m FavoriteModel) AddFavorite(userId, toolId int64) error {
	query := `
		INSERT INTO favorites (user_id, tool_id)
		VALUES ($1, $2)
		`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, userId, toolId)
	return err
}

func (m FavoriteModel) RemoveFavorite(userId, toolId int64) error {
	query := `
		DELETE FROM favorites
		WHERE user_id = $1 AND tool_id = $2
		`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, userId, toolId)
	return err
}

func (m FavoriteModel) GetFavorites(userId int64) ([]Favorite, error) {
	query := `
		SELECT user_id, tool_id
		FROM favorites
		WHERE user_id = $1
		`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var favorites []Favorite
	for rows.Next() {
		var favorite Favorite
		err := rows.Scan(
			&favorite.UserId,
			&favorite.ToolId,
		)
		if err != nil {
			return nil, err
		}
		favorites = append(favorites, favorite)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return favorites, nil
}
