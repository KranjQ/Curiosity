package psql

import (
	"context"
	"curiosity/internal/models"
	"database/sql"
	"fmt"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (repo *UserRepository) CreateUser(ctx context.Context, user models.User) (int, error) {
	query := `INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id`

	var id int
	row := repo.DB.QueryRowContext(ctx, query, user.Username, user.Password)
	if err := row.Scan(&id); err != nil {
		return -1, fmt.Errorf("bad scan user id: %w", err)
	}
	return id, nil
}

func (repo *UserRepository) GetUser(ctx context.Context, user models.User) (models.User, error) {
	query := "SELECT id, username, password FROM users WHERE username = $1"

	row := repo.DB.QueryRowContext(ctx, query, user.Username)

	var returnUser models.User

	if err := row.Scan(&returnUser.ID, &returnUser.Username, &returnUser.Password); err != nil {
		return models.User{}, fmt.Errorf("bad user scan error: %w", err)
	}
	return returnUser, nil
}
