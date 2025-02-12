package inmemory

import (
	"context"
	"curiosity/internal/models"
	"curiosity/internal/pkg/content/repo/inmemory/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCreateComment(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	cache := mocks.NewMockCommentCache(mockCtrl)

	repo := NewCommentRepository(cache)

	tests := []struct {
		name           string
		comment        models.Comment
		commentID      int
		cacheGetResult models.Comment
		cacheGetBool   bool
		expectedError  error
	}{
		{
			name: "Correct Create Comment",
			comment: models.Comment{
				Author:  1,
				Message: "message",
				Post:    1,
				Parent:  -1,
				Depth:   3,
				Path:    "1,2",
				Replies: 0,
			},
			commentID: 3,
			cacheGetResult: models.Comment{
				ID:        1,
				Author:    1,
				Message:   "message",
				Post:      1,
				Parent:    1,
				Depth:     1,
				Path:      "1",
				Replies:   0,
				CreatedAt: time.Date(2006, 2, 2, 2, 2, 2, 0, time.UTC),
			},
			cacheGetBool:  true,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache.EXPECT().Get(gomock.Any()).Return(tt.cacheGetResult, tt.cacheGetBool).Times(2)
			cache.EXPECT().Set(gomock.Any(), gomock.Any()).Return().Times(3)

			err := repo.CreateComment(ctx, tt.comment)
			require.ErrorIs(t, err, tt.expectedError)

		})
	}
}

func TestGetCommentByID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	cache := mocks.NewMockCommentCache(mockCtrl)

	repo := NewCommentRepository(cache)

	tests := []struct {
		name           string
		commentID      int
		cacheGetResult models.Comment
		cacheGetBool   bool
		expectedError  error
	}{
		{
			name:      "Correct Get Comment by ID",
			commentID: 1,
			cacheGetResult: models.Comment{
				ID:        1,
				Author:    1,
				Message:   "message",
				Post:      1,
				Parent:    -1,
				Depth:     1,
				Path:      "",
				Replies:   0,
				CreatedAt: time.Date(2006, 2, 2, 2, 2, 2, 0, time.UTC),
			},
			cacheGetBool:  true,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache.EXPECT().Get(tt.commentID).Return(tt.cacheGetResult, tt.cacheGetBool)

			result, err := repo.GetCommentByID(ctx, tt.commentID)
			require.ErrorIs(t, err, tt.expectedError)
			require.Equal(t, result, tt.cacheGetResult)
		})
	}
}

func TestGetCommentsByPostID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	cache := mocks.NewMockCommentCache(mockCtrl)

	repo := NewCommentRepository(cache)

	tests := []struct {
		name              string
		postID            int
		limit             int
		offset            int
		cacheGetAllResult []models.Comment
		expectedResult    []models.Comment
		expectedError     error
	}{
		{
			name:   "Correct Get Comment by Post ID",
			postID: 1,
			limit:  10,
			offset: 0,
			cacheGetAllResult: []models.Comment{
				{
					ID:        1,
					Author:    1,
					Message:   "message",
					Post:      1,
					Parent:    -1,
					Depth:     1,
					Path:      "",
					Replies:   0,
					CreatedAt: time.Date(2006, 2, 2, 2, 2, 2, 0, time.UTC),
				},
				{
					ID:        2,
					Author:    1,
					Message:   "message",
					Post:      1,
					Parent:    -1,
					Depth:     1,
					Path:      "",
					Replies:   0,
					CreatedAt: time.Date(2007, 2, 2, 2, 2, 2, 0, time.UTC),
				},
			},
			expectedResult: []models.Comment{
				{
					ID:        2,
					Author:    1,
					Message:   "message",
					Post:      1,
					Parent:    -1,
					Depth:     1,
					Path:      "",
					Replies:   0,
					CreatedAt: time.Date(2007, 2, 2, 2, 2, 2, 0, time.UTC),
				},
				{
					ID:        1,
					Author:    1,
					Message:   "message",
					Post:      1,
					Parent:    -1,
					Depth:     1,
					Path:      "",
					Replies:   0,
					CreatedAt: time.Date(2006, 2, 2, 2, 2, 2, 0, time.UTC),
				},
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache.EXPECT().GetAll().Return(tt.cacheGetAllResult)

			result, err := repo.GetCommentsByPostID(ctx, tt.postID, tt.limit, tt.offset)
			require.ErrorIs(t, err, tt.expectedError)
			require.Equal(t, result, tt.expectedResult)
		})
	}
}

func TestGetRepliesByCommentID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	cache := mocks.NewMockCommentCache(mockCtrl)

	repo := NewCommentRepository(cache)

	tests := []struct {
		name              string
		commentID         int
		cacheGetAllResult []models.Comment
		expectedReplies   []models.Comment
		expectedError     error
	}{
		{
			name:      "Correct Get Comment by ID",
			commentID: 1,
			cacheGetAllResult: []models.Comment{
				{
					ID:        2,
					Author:    1,
					Message:   "message",
					Post:      1,
					Parent:    1,
					Depth:     1,
					Path:      "1",
					Replies:   0,
					CreatedAt: time.Date(2006, 2, 2, 2, 2, 2, 0, time.UTC),
				},
				{
					ID:        3,
					Author:    1,
					Message:   "message",
					Post:      1,
					Parent:    1,
					Depth:     1,
					Path:      "1",
					Replies:   0,
					CreatedAt: time.Date(2007, 2, 2, 2, 2, 2, 0, time.UTC),
				},
			},
			expectedReplies: []models.Comment{
				{
					ID:        3,
					Author:    1,
					Message:   "message",
					Post:      1,
					Parent:    1,
					Depth:     1,
					Path:      "1",
					Replies:   0,
					CreatedAt: time.Date(2007, 2, 2, 2, 2, 2, 0, time.UTC),
				},
				{
					ID:        2,
					Author:    1,
					Message:   "message",
					Post:      1,
					Parent:    1,
					Depth:     1,
					Path:      "1",
					Replies:   0,
					CreatedAt: time.Date(2006, 2, 2, 2, 2, 2, 0, time.UTC),
				},
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache.EXPECT().GetAll().Return(tt.cacheGetAllResult)

			result, err := repo.GetRepliesByCommentID(ctx, tt.commentID)
			require.ErrorIs(t, err, tt.expectedError)
			require.Equal(t, result, tt.expectedReplies)
		})
	}
}
