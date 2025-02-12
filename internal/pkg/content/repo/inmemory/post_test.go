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

func TestCreatePost(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	cache := mocks.NewMockPostCache(mockCtrl)

	repo := NewPostRepository(cache)

	tests := []struct {
		name          string
		post          models.Post
		expectedError error
	}{
		{
			name: "Correct Create Post",
			post: models.Post{
				ID:            1,
				Title:         "Title",
				Content:       "Content",
				Author:        1,
				IsCommentable: false,
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache.EXPECT().Set(tt.post.ID, tt.post)

			err := repo.CreatePost(ctx, tt.post)
			require.ErrorIs(t, err, tt.expectedError)
		})
	}
}

func TestGetPosts(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	cache := mocks.NewMockPostCache(mockCtrl)

	repo := NewPostRepository(cache)

	tests := []struct {
		name              string
		cacheGetAllResult []models.Post
		expectedPosts     []models.Post
		expectedError     error
	}{
		{
			name: "Correct Get Posts",
			cacheGetAllResult: []models.Post{
				{
					ID:            1,
					Title:         "Title1",
					Content:       "Content1",
					Author:        1,
					IsCommentable: false,
					CreatedAt:     time.Date(2006, 2, 2, 2, 2, 2, 0, time.UTC),
				},
				{
					ID:            2,
					Title:         "Title2",
					Content:       "Content2",
					Author:        1,
					IsCommentable: false,
					CreatedAt:     time.Date(2007, 2, 2, 2, 2, 2, 0, time.UTC),
				},
			},
			expectedPosts: []models.Post{
				{
					ID:            2,
					Title:         "Title2",
					Content:       "Content2",
					Author:        1,
					IsCommentable: false,
					CreatedAt:     time.Date(2007, 2, 2, 2, 2, 2, 0, time.UTC),
				},
				{
					ID:            1,
					Title:         "Title1",
					Content:       "Content1",
					Author:        1,
					IsCommentable: false,
					CreatedAt:     time.Date(2006, 2, 2, 2, 2, 2, 0, time.UTC),
				},
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache.EXPECT().GetAll().Return(tt.cacheGetAllResult)

			result, err := repo.GetPosts(ctx)
			require.ErrorIs(t, err, tt.expectedError)
			require.Equal(t, result, tt.expectedPosts)
		})
	}
}

func TestGetPostByID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	cache := mocks.NewMockPostCache(mockCtrl)

	repo := NewPostRepository(cache)

	tests := []struct {
		name           string
		postID         int
		cacheGetResult models.Post
		cacheGetExists bool
		expectedPost   models.Post
		expectedError  error
	}{
		{
			name:   "Correct Get Post By ID",
			postID: 1,
			cacheGetResult: models.Post{
				ID:            1,
				Title:         "Title",
				Content:       "Content",
				Author:        1,
				IsCommentable: false,
				CreatedAt:     time.Date(2006, 2, 2, 2, 2, 2, 0, time.UTC),
			},
			cacheGetExists: true,
			expectedPost: models.Post{
				ID:            1,
				Title:         "Title",
				Content:       "Content",
				Author:        1,
				IsCommentable: false,
				CreatedAt:     time.Date(2006, 2, 2, 2, 2, 2, 0, time.UTC),
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		cache.EXPECT().Get(tt.postID).Return(tt.cacheGetResult, tt.cacheGetExists)

		result, err := repo.GetPostByID(ctx, tt.postID)
		require.ErrorIs(t, err, tt.expectedError)
		require.Equal(t, result, tt.expectedPost)
	}
}
