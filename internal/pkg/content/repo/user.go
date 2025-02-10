package repo

import (
	"context"
	"curiosity/internal/models"
	"database/sql"
	"fmt"
	"log"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (repo *UserRepository) CreateUser(ctx context.Context, user models.User) (*models.User, error) {
	query := `INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id`

	var returnUser models.User
	log.Print("start repo")
	row := repo.DB.QueryRowContext(ctx, query, user.Username, user.Password)
	if err := row.Scan(&returnUser.ID); err != nil {
		return nil, fmt.Errorf("bad insert user in db: %w", err)
	}
	log.Print("end repo")
	return &returnUser, nil
}

func (repo *UserRepository) GetUser(ctx context.Context, user models.User) (*models.User, error) {
	query := "SELECT id, username, password FROM users WHERE username = $1"

	row := repo.DB.QueryRowContext(ctx, query, user.Username)

	var returnUser models.User

	if err := row.Scan(&returnUser.ID, &returnUser.Username, &returnUser.Password); err != nil {
		return nil, fmt.Errorf("bad user scan error: %w", err)
	}
	return &returnUser, nil
}
