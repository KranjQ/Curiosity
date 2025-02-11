package usecase

import (
	"context"
	"curiosity/internal/models"
	"fmt"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user models.User) (int, error)
	GetUser(ctx context.Context, user models.User) (models.User, error)
}

type UserUseCase struct {
	repo UserRepository
}

func NewUserUseCase(repo UserRepository) *UserUseCase {
	return &UserUseCase{repo: repo}
}

func (uc *UserUseCase) RegisterUser(ctx context.Context, user models.User) (int, error) {
	id, err := uc.repo.CreateUser(ctx, user)
	if err != nil {
		return -1, fmt.Errorf("create user error: %w", err)
	}
	return id, nil
}

func (uc *UserUseCase) SignIn(ctx context.Context, user models.User) (int, error) {
	u, err := uc.repo.GetUser(ctx, user)
	if err != nil {
		return -1, fmt.Errorf("bad get password: %w", err)
	}
	if user.Password != u.Password {
		return -1, fmt.Errorf("incorrect password error for: %s", user.Username)
	}
	return u.ID, nil
}
