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
				Title:         "Тестовый пост",
				Content:       "Это тестовый пост",
				Author:        1,
				IsCommentable: false,
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectExec("INSERT INTO posts").
				WillReturnResult(sqlmock.NewResult(1, 1)).
				WithArgs(tt.post.Title, tt.post.Content, tt.post.Author, tt.post.IsCommentable).
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
		wantResult []models.Post
	}{
		{
			name: "Correct Get Posts",
			dbResult: sqlmock.NewRows([]string{"id", "title", "content", "author", "is_commentable", "created_at"}).
				AddRow(1, "Тестовый заголовок", "Тестовое содержание", 1, false, time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)),
			dbError: nil,
			wantResult: []models.Post{{ID: 1, Title: "Тестовый заголовок", Content: "Тестовое содержание", Author: 1,
				IsCommentable: false, CreatedAt: time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)}},
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
			mock.ExpectQuery("SELECT id, title, content, author, is_commentable, created_at FROM posts").
				WillReturnRows(tt.dbResult).
				WillReturnError(tt.dbError)
			posts, err := repo.GetPosts(ctx)
			require.ErrorIs(t, err, tt.dbError)
			require.Equal(t, posts, tt.wantResult)
		})
	}
}

func TestGetPostByID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sql new mock error: %v", err)
	}
	defer db.Close()

	repo := NewPostRepository(db)

	tests := []struct {
		name       string
		postID     int
		dbRow      *sqlmock.Rows
		dbError    error
		expectPost models.Post
	}{
		{
			name:   "Correct select",
			postID: 1,
			dbRow: sqlmock.NewRows([]string{"id", "title", "content", "author", "is_commentable", "created_at"}).
				AddRow(1, "Тестовый тайтл", "Тестовый контент", 1, false, time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)),
			dbError: nil,
			expectPost: models.Post{
				ID:            1,
				Title:         "Тестовый тайтл",
				Content:       "Тестовый контент",
				Author:        1,
				IsCommentable: false,
				CreatedAt:     time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectQuery("SELECT id, title, content, author, is_commentable, created_at  FROM posts").
				WithArgs(tt.postID).
				WillReturnRows(tt.dbRow).
				WillReturnError(tt.dbError)

			result, err := repo.GetPostByID(ctx, tt.postID)
			require.ErrorIs(t, err, tt.dbError)
			require.Equal(t, result, tt.expectPost)
		})
	}

}
