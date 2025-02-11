package usecase

import (
	"context"
	"curiosity/internal/models"
	"fmt"
	"strconv"
)

//go:generate mockgen -destination=./mocks/mock_commentRepository.go -package=mocks . CommentRepository
type CommentRepository interface {
	CreateComment(ctx context.Context, comment models.Comment) error
	GetCommentByID(ctx context.Context, id int) (models.Comment, error)
	GetCommentsByPostID(ctx context.Context, postID int, limit int, offset int) ([]models.Comment, error)
	GetRepliesByCommentID(ctx context.Context, commentID int) ([]models.Comment, error)
}

//go:generate mockgen -destination=./mocks/mock_cPostRepository.go -package=mocks . CPostRepository
type CPostRepository interface {
	GetPostByID(ctx context.Context, postID int) (models.Post, error)
}

type CommentUseCase struct {
	commentRepo CommentRepository
	postRepo    CPostRepository
}

func NewCommentUseCase(cr CommentRepository, pr CPostRepository) *CommentUseCase {
	return &CommentUseCase{
		commentRepo: cr,
		postRepo:    pr,
	}
}

func (uc *CommentUseCase) CreateComment(ctx context.Context, comment models.Comment) error {
	post, err := uc.postRepo.GetPostByID(ctx, comment.Post)
	if err != nil {
		return fmt.Errorf("bad GetPostByID error: %w", err)
	}
	if !post.IsCommentable {
		return fmt.Errorf("Post is not commentable, PostID: %d", post.ID)
	}
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

	if err := uc.commentRepo.CreateComment(ctx, comment); err != nil {
		return fmt.Errorf("CreateComment error: %w", err)
	}
	return nil
}

func (uc *CommentUseCase) GetCommentByID(ctx context.Context, id int) (models.Comment, error) {
	comment, err := uc.commentRepo.GetCommentByID(ctx, id)
	if err != nil {
		return models.Comment{}, fmt.Errorf("GetCommentByID error: %w", err)
	}
	return comment, nil
}

func (uc *CommentUseCase) GetCommentsByPostID(ctx context.Context, postID int, limit int, offset int) ([]models.Comment, error) {
	comments, err := uc.commentRepo.GetCommentsByPostID(ctx, postID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("bad repo GetCommentsByPostID: %w", err)
	}
	return comments, nil
}

func (uc *CommentUseCase) GetRepliesByCommentID(ctx context.Context, commentID int) ([]models.Comment, error) {
	comments, err := uc.commentRepo.GetRepliesByCommentID(ctx, commentID)
	if err != nil {
		return nil, fmt.Errorf("bad repo GetCommentsByPostID: %w", err)
	}
	return comments, nil
}
