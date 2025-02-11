package usecase

import (
	"context"
	"curiosity/internal/models"
	"curiosity/internal/pkg/content/usecase/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestRegisterUser(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	userRepo := mocks.NewMockUserRepository(mockCtrl)

	uc := NewUserUseCase(userRepo)

	tests := []struct {
		name             string
		user             models.User
		createUserResult int
		createUserError  error
		expectedID       int
	}{
		{
			name: "Correct Register User",
			user: models.User{
				Username: "Alice",
				Password: "123",
			},
			createUserResult: 1,
			createUserError:  nil,
			expectedID:       1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo.EXPECT().CreateUser(ctx, tt.user).Return(tt.createUserResult, tt.createUserError)

			result, err := uc.RegisterUser(ctx, tt.user)
			require.ErrorIs(t, err, tt.createUserError)
			require.Equal(t, result, tt.createUserResult)
		})
	}
}

func TestSignIn(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	userRepo := mocks.NewMockUserRepository(mockCtrl)

	uc := NewUserUseCase(userRepo)

	tests := []struct {
		name          string
		user          models.User
		getUserResult models.User
		getUserError  error
		expectedID    int
	}{
		{
			name: "Correct Sign In",
			user: models.User{
				Username: "Alice",
				Password: "123",
			},
			getUserResult: models.User{
				ID:       1,
				Username: "Alice",
				Password: "123",
			},
			getUserError: nil,
			expectedID:   1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo.EXPECT().GetUser(ctx, tt.user).Return(tt.getUserResult, tt.getUserError)

			result, err := uc.SignIn(ctx, tt.user)

			require.ErrorIs(t, err, tt.getUserError)
			require.Equal(t, result, tt.getUserResult.ID)
		})
	}
}
