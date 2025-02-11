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

func TestCreateComment(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	commentRepository := mocks.NewMockCommentRepository(mockCtrl)
	cPostRepository := mocks.NewMockCPostRepository(mockCtrl)

	uc := NewCommentUseCase(commentRepository, cPostRepository)

	tests := []struct {
		name                   string
		comment                models.Comment
		getPostByIDResult      models.Post
		getPostByIDError       error
		getCommentByIDResult   models.Comment
		getCommentByIDError    error
		createCommentParameter models.Comment
		createCommentError     error
		expectedError          error
	}{
		{
			name: "Correct Create Comment",
			comment: models.Comment{
				Author:  1,
				Message: "message",
				Post:    1,
				Parent:  1,
			},
			getPostByIDResult: models.Post{
				IsCommentable: true,
			},
			getPostByIDError: nil,
			getCommentByIDResult: models.Comment{
				ID:     1,
				Parent: -1,
				Depth:  1,
				Path:   "",
			},
			getCommentByIDError: nil,
			createCommentParameter: models.Comment{
				Author:  1,
				Message: "message",
				Post:    1,
				Parent:  1,
				Depth:   2,
				Path:    "1",
			},
			createCommentError: nil,
			expectedError:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cPostRepository.EXPECT().GetPostByID(ctx, tt.comment.Post).Return(tt.getPostByIDResult, tt.getPostByIDError)
			commentRepository.EXPECT().GetCommentByID(ctx, tt.comment.Parent).Return(tt.getCommentByIDResult, tt.getCommentByIDError)
			commentRepository.EXPECT().CreateComment(ctx, tt.createCommentParameter).Return(tt.createCommentError)

			result := uc.CreateComment(ctx, tt.comment)
			require.ErrorIs(t, result, tt.expectedError)
		})
	}
}

func TestGetCommentByID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	commentRepo := mocks.NewMockCommentRepository(mockCtrl)
	cPostRepo := mocks.NewMockCPostRepository(mockCtrl)

	uc := NewCommentUseCase(commentRepo, cPostRepo)

	tests := []struct {
		name                 string
		id                   int
		getCommentByIDResult models.Comment
		getCommentByIDError  error
		expectedComment      models.Comment
	}{
		{
			name: "Correct Get Comment by ID",
			id:   1,
			getCommentByIDResult: models.Comment{
				ID:        1,
				Author:    1,
				Message:   "message",
				Post:      1,
				Parent:    -1,
				Depth:     1,
				Path:      "",
				Replies:   0,
				CreatedAt: time.Date(2006, 1, 2, 15, 2, 3, 0, time.UTC),
			},
			getCommentByIDError: nil,
			expectedComment: models.Comment{
				ID:        1,
				Author:    1,
				Message:   "message",
				Post:      1,
				Parent:    -1,
				Depth:     1,
				Path:      "",
				Replies:   0,
				CreatedAt: time.Date(2006, 1, 2, 15, 2, 3, 0, time.UTC),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			commentRepo.EXPECT().GetCommentByID(ctx, tt.id).Return(tt.getCommentByIDResult, tt.getCommentByIDError)

			result, err := uc.GetCommentByID(ctx, tt.id)

			require.ErrorIs(t, err, tt.getCommentByIDError)
			require.Equal(t, result, tt.expectedComment)
		})
	}
}

func TestGetCommentsByPostID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	commentRepo := mocks.NewMockCommentRepository(mockCtrl)
	cPostRepo := mocks.NewMockCPostRepository(mockCtrl)

	uc := NewCommentUseCase(commentRepo, cPostRepo)

	tests := []struct {
		name                      string
		postID                    int
		limit                     int
		offset                    int
		getCommentsByPostIDResult []models.Comment
		getCommentsByPostIDError  error
		expectedComments          []models.Comment
	}{
		{
			name:   "Correct Get Comment by ID",
			postID: 1,
			limit:  10,
			offset: 0,
			getCommentsByPostIDResult: []models.Comment{
				{
					ID:        1,
					Author:    1,
					Message:   "message",
					Post:      1,
					Parent:    -1,
					Depth:     1,
					Path:      "",
					Replies:   0,
					CreatedAt: time.Date(2006, 1, 2, 15, 2, 3, 0, time.UTC),
				},
			},
			getCommentsByPostIDError: nil,
			expectedComments: []models.Comment{
				{
					ID:        1,
					Author:    1,
					Message:   "message",
					Post:      1,
					Parent:    -1,
					Depth:     1,
					Path:      "",
					Replies:   0,
					CreatedAt: time.Date(2006, 1, 2, 15, 2, 3, 0, time.UTC),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			commentRepo.EXPECT().GetCommentsByPostID(ctx, tt.postID, tt.limit, tt.offset).Return(tt.getCommentsByPostIDResult, tt.getCommentsByPostIDError)

			result, err := uc.GetCommentsByPostID(ctx, tt.postID, tt.limit, tt.offset)
			require.ErrorIs(t, err, tt.getCommentsByPostIDError)
			require.Equal(t, result, tt.getCommentsByPostIDResult)
		})
	}
}

func TestGetRepliesByCommentID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	commentRepo := mocks.NewMockCommentRepository(mockCtrl)
	cPostRepo := mocks.NewMockCPostRepository(mockCtrl)

	uc := NewCommentUseCase(commentRepo, cPostRepo)

	tests := []struct {
		name                        string
		commentID                   int
		getRepliesByCommentIDResult []models.Comment
		getRepliesByCommentIDError  error
		expectedComments            []models.Comment
	}{
		{
			name:      "Correct Get Comment by ID",
			commentID: 1,
			getRepliesByCommentIDResult: []models.Comment{
				{
					ID:        3,
					Author:    1,
					Message:   "message",
					Post:      1,
					Parent:    2,
					Depth:     2,
					Path:      "",
					Replies:   0,
					CreatedAt: time.Date(2006, 1, 2, 15, 2, 3, 0, time.UTC),
				},
			},
			getRepliesByCommentIDError: nil,
			expectedComments: []models.Comment{
				{
					ID:        3,
					Author:    1,
					Message:   "message",
					Post:      1,
					Parent:    2,
					Depth:     2,
					Path:      "",
					Replies:   0,
					CreatedAt: time.Date(2006, 1, 2, 15, 2, 3, 0, time.UTC),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			commentRepo.EXPECT().GetRepliesByCommentID(ctx, tt.commentID).Return(tt.getRepliesByCommentIDResult, tt.getRepliesByCommentIDError)

			result, err := uc.GetRepliesByCommentID(ctx, tt.commentID)
			require.ErrorIs(t, err, tt.getRepliesByCommentIDError)
			require.Equal(t, result, tt.getRepliesByCommentIDResult)
		})
	}
}
