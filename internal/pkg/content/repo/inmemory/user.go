package inmemory

import (
	"context"
	"curiosity/internal/models"
	"fmt"
	"sync"
)

type UserRepository struct {
	users   map[string]models.User
	mu      *sync.RWMutex
	counter int
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		users:   make(map[string]models.User),
		mu:      &sync.RWMutex{},
		counter: 1,
	}
}

func (repo *UserRepository) CreateUser(ctx context.Context, user models.User) (int, error) {
	_ = ctx
	repo.mu.Lock()
	defer repo.mu.Unlock()
	user.ID = repo.counter
	if _, exists := repo.users[user.Username]; exists {
		return -1, fmt.Errorf("This username is taken: %s", user.Username)
	}
	repo.users[user.Username] = user
	repo.counter += 1
	return user.ID, nil
}

func (repo *UserRepository) GetUser(ctx context.Context, user models.User) (models.User, error) {
	_ = ctx
	repo.mu.RLock()
	defer repo.mu.RUnlock()
	returnUser, exists := repo.users[user.Username]
	if !exists {
		return models.User{}, fmt.Errorf("No user with that username")
	}
	return returnUser, nil
}
