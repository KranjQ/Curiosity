package psql

import (
	"context"
	"curiosity/internal/models"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCreatePost(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("bad sql mock new: %v", err)
	}
	defer db.Close()

	repo := NewPostRepository(db)

	tests := []struct {
		name    string
		post    models.Post
		wantErr error
	}{
		{
			name: "Correct Post Data",
			post: models.Post{
				Title:   "Тестовый пост",
				Content: "Это тестовый пост",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectExec("INSERT INTO posts").
				WillReturnResult(sqlmock.NewResult(1, 1)).
				WithArgs(tt.post.Title, tt.post.Content).
				WillReturnError(nil)

			err := repo.CreatePost(ctx, tt.post)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestGetPosts(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("bad sqlmock db: %v", err)
	}
	defer db.Close()

	repo := NewPostRepository(db)

	tests := []struct {
		name       string
		dbResult   *sqlmock.Rows
		dbError    error
		wantResult []*models.Post
	}{
		{
			name: "Correct Get Posts",
			dbResult: sqlmock.NewRows([]string{"title", "content"}).
				AddRow("Тестовый заголовок", "Тестовое содержание"),
			dbError:    nil,
			wantResult: []*models.Post{{Title: "Тестовый заголовок", Content: "Тестовое содержание"}},
		},
		{
			name:       "DB error",
			dbResult:   sqlmock.NewRows([]string{}),
			dbError:    errors.New("DB error"),
			wantResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectQuery("SELECT id, title, content FROM posts").
				WillReturnRows(tt.dbResult).
				WillReturnError(tt.dbError)
			posts, err := repo.GetPosts(ctx)
			require.ErrorIs(t, err, tt.dbError)
			require.Equal(t, posts, tt.wantResult)
		})
	}
}

//func TestPostByID(t *testing.T) {
//	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
//	defer cancel()
//
//}
