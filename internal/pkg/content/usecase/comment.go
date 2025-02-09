package usecase

import (
	"context"
	"curiosity/internal/models"
	"fmt"
	"strconv"
)

type CommentRepository interface {
	CreateComment(ctx context.Context, comment models.Comment) error
	GetCommentByID(ctx context.Context, id int) (*models.Comment, error)
	GetCommentsByPostID(ctx context.Context, postID int, limit int, offset int) ([]*models.Comment, error)
	GetRepliesByCommentID(ctx context.Context, commentID int) ([]*models.Comment, error)
}

type CommentUseCase struct {
	repo CommentRepository
}

func NewCommentUseCase(repo CommentRepository) *CommentUseCase {
	return &CommentUseCase{repo: repo}
}

func (uc *CommentUseCase) CreateComment(ctx context.Context, comment models.Comment) error {
	if comment.Parent != -1 {
		parent, err := uc.GetCommentByID(ctx, comment.Parent)
		if err != nil {
			return fmt.Errorf("CreateComment GetParent error: %w", err)
		}
		comment.Depth = parent.Depth + 1
		if parent.Path == "" {
			comment.Path = strconv.Itoa(parent.ID)
		} else {
			comment.Path = parent.Path + "," + strconv.Itoa(parent.ID)
		}
	} else {
		comment.Path = ""
		comment.Depth = 1
	}

	comment.Author = 1

	if err := uc.repo.CreateComment(ctx, comment); err != nil {
		return fmt.Errorf("CreateComment error: %w", err)
	}
	return nil
}

func (uc *CommentUseCase) GetCommentByID(ctx context.Context, id int) (*models.Comment, error) {
	comment, err := uc.repo.GetCommentByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("GetCommentByID error: %w", err)
	}
	return comment, nil
}

func (uc *CommentUseCase) GetCommentsByPostID(ctx context.Context, postID int, limit int, offset int) ([]*models.Comment, error) {
	comments, err := uc.repo.GetCommentsByPostID(ctx, postID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("bad repo GetCommentsByPostID: %w", err)
	}
	return comments, nil
}

func (uc *CommentUseCase) GetRepliesByCommentID(ctx context.Context, commentID int) ([]*models.Comment, error) {
	comments, err := uc.repo.GetRepliesByCommentID(ctx, commentID)
	if err != nil {
		return nil, fmt.Errorf("bad repo GetCommentsByPostID: %w", err)
	}
	return comments, nil
}
