package psql

import (
	"context"
	"curiosity/internal/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCreateComment(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock new error: %v", err)
	}
	repo := NewCommentRepository(db)

	tests := []struct {
		name          string
		comment       models.Comment
		dbInsertRow   *sqlmock.Rows
		dbInsertError error
		expectedError error
	}{
		{
			name: "Correct create comment",
			comment: models.Comment{
				ID:      1,
				Author:  1,
				Message: "message",
				Post:    1,
				Parent:  -1,
				Depth:   1,
				Path:    "",
				Replies: 0,
			},
			dbInsertRow:   sqlmock.NewRows([]string{"id"}).AddRow(1),
			dbInsertError: nil,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectBegin()
			mock.ExpectQuery("INSERT INTO").
				WithArgs(tt.comment.Author, tt.comment.Message,
					tt.comment.Post, tt.comment.Parent, tt.comment.Depth, tt.comment.Path).
				WillReturnRows(tt.dbInsertRow).
				WillReturnError(tt.dbInsertError)
			mock.ExpectCommit()
			err := repo.CreateComment(ctx, tt.comment)
			require.ErrorIs(t, err, tt.expectedError)
		})
	}
}

func TestGetCommentByID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock new error: %v", err)
	}

	repo := NewCommentRepository(db)

	tests := []struct {
		name            string
		id              int
		dbRow           *sqlmock.Rows
		dbError         error
		expectedComment models.Comment
		expectedError   error
	}{
		{
			name: "Correct Get Comment by ID",
			id:   1,
			dbRow: sqlmock.NewRows([]string{"id", "author", "message", "post", "parent", "depth", "path", "replies", "created_at"}).
				AddRow(1, 1, "message", 1, -1, 1, "", 0, time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)),
			dbError: nil,
			expectedComment: models.Comment{
				ID:        1,
				Author:    1,
				Message:   "message",
				Post:      1,
				Parent:    -1,
				Depth:     1,
				Path:      "",
				Replies:   0,
				CreatedAt: time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC),
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectQuery("SELECT id, author, message, post, parent, depth, path, replies, created_at FROM comments").
				WithArgs(tt.id).
				WillReturnRows(tt.dbRow).
				WillReturnError(tt.dbError)

			result, err := repo.GetCommentByID(ctx, tt.id)
			require.ErrorIs(t, err, tt.expectedError)
			require.Equal(t, result, tt.expectedComment)
		})
	}
}

func TestGetCommentsByPostID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock new error: %v", err)
	}

	repo := NewCommentRepository(db)

	tests := []struct {
		name             string
		postID           int
		limit            int
		offset           int
		dbRow            *sqlmock.Rows
		dbError          error
		expectedComments []models.Comment
		expectedError    error
	}{
		{
			name:   "Correct Get Comments by Post ID",
			postID: 1,
			limit:  10,
			offset: 0,
			dbRow: sqlmock.NewRows([]string{"id", "author", "message", "post", "parent", "depth", "path", "replies", "created_at"}).
				AddRow(1, 1, "message", 1, -1, 1, "", 0, time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)),
			dbError: nil,
			expectedComments: []models.Comment{
				models.Comment{
					ID:        1,
					Author:    1,
					Message:   "message",
					Post:      1,
					Parent:    -1,
					Depth:     1,
					Path:      "",
					Replies:   0,
					CreatedAt: time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC),
				},
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectQuery("SELECT id, author, message, post, parent, depth, path, replies, created_at FROM comments").
				WithArgs(tt.postID, tt.limit, tt.offset).
				WillReturnRows(tt.dbRow).
				WillReturnError(tt.dbError)

			result, err := repo.GetCommentsByPostID(ctx, tt.postID, tt.limit, tt.offset)
			require.ErrorIs(t, err, tt.expectedError)
			require.Equal(t, result, tt.expectedComments)
		})
	}
}

func TestGetRepliesByCommentID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock new error: %v", err)
	}

	repo := NewCommentRepository(db)

	tests := []struct {
		name             string
		commentID        int
		dbRow            *sqlmock.Rows
		dbError          error
		expectedComments []models.Comment
		expectedError    error
	}{
		{
			name:      "Correct Get Comment by ID",
			commentID: 1,
			dbRow: sqlmock.NewRows([]string{"id", "author", "message", "post", "parent", "depth", "path", "replies", "created_at"}).
				AddRow(1, 1, "message", 1, -1, 1, "", 0, time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)),
			dbError: nil,
			expectedComments: []models.Comment{
				models.Comment{
					ID:        1,
					Author:    1,
					Message:   "message",
					Post:      1,
					Parent:    -1,
					Depth:     1,
					Path:      "",
					Replies:   0,
					CreatedAt: time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC),
				},
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectQuery("SELECT id, author, message, post, parent, depth, path, replies, created_at FROM comments").
				WithArgs(tt.commentID).
				WillReturnRows(tt.dbRow).
				WillReturnError(tt.dbError)

			result, err := repo.GetRepliesByCommentID(ctx, tt.commentID)
			require.ErrorIs(t, err, tt.expectedError)
			require.Equal(t, result, tt.expectedComments)
		})
	}
}
