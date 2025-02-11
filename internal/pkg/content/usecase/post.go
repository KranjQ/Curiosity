package usecase

import (
	"context"
	"curiosity/internal/models"
	"fmt"
)

type PostRepository interface {
	CreatePost(ctx context.Context, post models.Post) error
	GetPostByID(ctx context.Context, postID int) (models.Post, error)
	GetPosts(ctx context.Context) ([]models.Post, error)
}

type PostUseCase struct {
	repo PostRepository
}

func NewPostUseCase(repo PostRepository) *PostUseCase {
	return &PostUseCase{repo: repo}
}

func (uc *PostUseCase) CreatePost(ctx context.Context, post models.Post) error {
	if err := uc.repo.CreatePost(ctx, post); err != nil {
		return fmt.Errorf("CreatePost error: %w", err)
	}
	return nil
}

func (uc *PostUseCase) GetPostByID(ctx context.Context, postID int) (models.Post, error) {
	post, err := uc.repo.GetPostByID(ctx, postID)
	if err != nil {
		return models.Post{}, fmt.Errorf("bad GetPostByID error: %w", err)
	}
	return post, nil
}

func (uc *PostUseCase) GetPosts(ctx context.Context) ([]models.Post, error) {
	posts, err := uc.repo.GetPosts(ctx)
	if err != nil {
		return nil, fmt.Errorf("bad GetPosts error: %w", err)
	}
	return posts, nil
}
