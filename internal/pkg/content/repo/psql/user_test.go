package psql

import (
	"context"
	"curiosity/internal/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCreateUser(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("slqmokc New error: %v", err)
	}

	repo := NewUserRepository(db)

	tests := []struct {
		name       string
		user       models.User
		dbRow      *sqlmock.Rows
		dbError    error
		expectedID int
	}{
		{
			name: "Correct create user",
			user: models.User{
				Username: "Alice",
				Password: "123",
			},
			dbRow:      sqlmock.NewRows([]string{"id"}).AddRow(1),
			dbError:    nil,
			expectedID: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectQuery("INSERT INTO users").
				WithArgs(tt.user.Username, tt.user.Password).
				WillReturnRows(tt.dbRow).
				WillReturnError(tt.dbError)

			id, err := repo.CreateUser(ctx, tt.user)
			require.ErrorIs(t, err, tt.dbError)
			require.Equal(t, id, tt.expectedID)
		})
	}
}

func TestGetUser(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock new error: %v", err)
	}

	repo := NewUserRepository(db)

	tests := []struct {
		name         string
		user         models.User
		dbRows       *sqlmock.Rows
		dbError      error
		expectedUser models.User
	}{
		{
			name: "Coreect Get User",
			user: models.User{Username: "Alice"},
			dbRows: sqlmock.NewRows([]string{"id", "username", "password"}).
				AddRow(1, "Alice", "123"),
			dbError: nil,
			expectedUser: models.User{
				ID:       1,
				Username: "Alice",
				Password: "123",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectQuery("SELECT id, username, password FROM users").
				WithArgs(tt.user.Username).
				WillReturnRows(tt.dbRows).
				WillReturnError(tt.dbError)

			result, err := repo.GetUser(ctx, tt.user)
			require.ErrorIs(t, err, tt.dbError)
			require.Equal(t, result, tt.expectedUser)
		})
	}
}
