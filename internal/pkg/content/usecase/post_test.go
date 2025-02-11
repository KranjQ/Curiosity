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

func TestCreatePost(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	postRepo := mocks.NewMockPostRepository(mockCtrl)

	uc := NewPostUseCase(postRepo)

	tests := []struct {
		name            string
		post            models.Post
		createPostError error
	}{
		{
			name: "Correct Create Post",
			post: models.Post{
				Title:         "Тестовой заголовок",
				Content:       "Тестовой контент",
				Author:        1,
				IsCommentable: false,
			},
			createPostError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			postRepo.EXPECT().CreatePost(ctx, tt.post).Return(tt.createPostError)

			err := uc.CreatePost(ctx, tt.post)
			require.ErrorIs(t, err, tt.createPostError)
		})
	}
}

func TestGetPostByID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	postRepo := mocks.NewMockPostRepository(mockCtrl)

	uc := NewPostUseCase(postRepo)

	tests := []struct {
		name              string
		postID            int
		getPostByIDResult models.Post
		getPostByIDError  error
		expectedPost      models.Post
	}{
		{
			name:   "Correct Get Post By ID",
			postID: 1,
			getPostByIDResult: models.Post{
				ID:            1,
				Title:         "Тестовый заголовок",
				Content:       "Тестовый контент",
				Author:        1,
				IsCommentable: false,
				CreatedAt:     time.Date(2006, 2, 2, 2, 2, 2, 0, time.UTC),
			},
			getPostByIDError: nil,
			expectedPost: models.Post{
				ID:            1,
				Title:         "Тестовый заголовок",
				Content:       "Тестовый контент",
				Author:        1,
				IsCommentable: false,
				CreatedAt:     time.Date(2006, 2, 2, 2, 2, 2, 0, time.UTC),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			postRepo.EXPECT().GetPostByID(ctx, tt.postID).Return(tt.getPostByIDResult, tt.getPostByIDError)

			result, err := uc.GetPostByID(ctx, tt.postID)
			require.ErrorIs(t, err, tt.getPostByIDError)
			require.Equal(t, result, tt.expectedPost)
		})
	}
}

func TestGetPosts(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	postRepo := mocks.NewMockPostRepository(mockCtrl)

	uc := NewPostUseCase(postRepo)

	tests := []struct {
		name           string
		getPostsResult []models.Post
		getPostsError  error
		expectedPosts  []models.Post
	}{
		{
			name: "Correct Get Post By ID",
			getPostsResult: []models.Post{
				{
					ID:            1,
					Title:         "Тестовый заголовок",
					Content:       "Тестовый контент",
					Author:        1,
					IsCommentable: false,
					CreatedAt:     time.Date(2006, 2, 2, 2, 2, 2, 0, time.UTC),
				},
			},
			getPostsError: nil,
			expectedPosts: []models.Post{
				{
					ID:            1,
					Title:         "Тестовый заголовок",
					Content:       "Тестовый контент",
					Author:        1,
					IsCommentable: false,
					CreatedAt:     time.Date(2006, 2, 2, 2, 2, 2, 0, time.UTC),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			postRepo.EXPECT().GetPosts(ctx).Return(tt.getPostsResult, tt.getPostsError)

			result, err := uc.GetPosts(ctx)
			require.ErrorIs(t, err, tt.getPostsError)
			require.Equal(t, result, tt.getPostsResult)
		})
	}
}
