package inmemory

import (
	"context"
	"curiosity/internal/models"
	"curiosity/internal/pkg/content/repo/inmemory/mocks"
	"curiosity/internal/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCreateUser(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	cache := mocks.NewMockUserCache(mockCtrl)

	repo := NewUserRepository(cache)

	tests := []struct {
		name           string
		user           models.User
		cacheGetResult models.User
		cacheGetExists bool
		expectedID     int
		expectedError  error
	}{
		{
			name: "Correct Create User",
			user: models.User{
				ID:       1,
				Username: "Alice",
				Password: "123",
			},
			cacheGetResult: models.User{},
			cacheGetExists: false,
			expectedID:     1,
			expectedError:  nil,
		},
		{
			name: "Username is taken",
			user: models.User{
				ID:       1,
				Username: "Alice",
				Password: "123",
			},
			cacheGetResult: models.User{
				ID:       1,
				Username: "Alice",
				Password: "123",
			},
			cacheGetExists: true,
			expectedID:     -1,
			expectedError:  utils.ErrUsernameExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectedError != utils.ErrUsernameExists {
				cache.EXPECT().Set(tt.user.Username, tt.user).Return().Times(1)
			}
			cache.EXPECT().Get(tt.user.Username).Return(tt.cacheGetResult, tt.cacheGetExists)

			result, err := repo.CreateUser(ctx, tt.user)
			require.ErrorIs(t, err, tt.expectedError)
			require.Equal(t, result, tt.expectedID)
		})
	}
}

func TestGetUser(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	cache := mocks.NewMockUserCache(mockCtrl)

	repo := NewUserRepository(cache)

	tests := []struct {
		name           string
		user           models.User
		cacheGetResult models.User
		cacheGetExists bool
		expectedUser   models.User
		expectedError  error
	}{
		{
			name: "Correct Get User",
			user: models.User{
				Username: "Alice",
				Password: "123",
			},
			cacheGetResult: models.User{
				ID:       1,
				Username: "Alice",
				Password: "123",
			},
			cacheGetExists: true,
			expectedUser: models.User{
				ID:       1,
				Username: "Alice",
				Password: "123",
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache.EXPECT().Get(tt.user.Username).Return(tt.cacheGetResult, tt.cacheGetExists)

			result, err := repo.GetUser(ctx, tt.user)
			require.ErrorIs(t, err, tt.expectedError)
			require.Equal(t, result, tt.cacheGetResult)
		})
	}
}
