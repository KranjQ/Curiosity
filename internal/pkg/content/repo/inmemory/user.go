package inmemory

import (
	"context"
	"curiosity/internal/models"
	"curiosity/internal/utils"
	"fmt"
	"sync"
)

//go:generate mockgen -destination=./mocks/mock_userCache.go -package=mocks . UserCache
type UserCache interface {
	Get(key string) (models.User, bool)
	Set(key string, value models.User)
	GetAll() []models.User
}

type UserRepository struct {
	cache   UserCache
	mu      *sync.RWMutex
	counter int
}

func NewUserRepository(uc UserCache) *UserRepository {
	return &UserRepository{
		cache:   uc,
		mu:      &sync.RWMutex{},
		counter: 0,
	}
}

func (repo *UserRepository) CreateUser(ctx context.Context, user models.User) (int, error) {
	_ = ctx
	repo.mu.Lock()
	repo.counter += 1
	user.ID = repo.counter
	repo.mu.Unlock()
	_, exists := repo.cache.Get(user.Username)
	if exists {
		return -1, utils.ErrUsernameExists
	}
	repo.cache.Set(user.Username, user)

	return user.ID, nil
}

func (repo *UserRepository) GetUser(ctx context.Context, user models.User) (models.User, error) {
	_ = ctx
	repo.mu.RLock()
	defer repo.mu.RUnlock()
	returnUser, exists := repo.cache.Get(user.Username)
	if !exists {
		return models.User{}, fmt.Errorf("No user with that username")
	}
	return returnUser, nil
}
