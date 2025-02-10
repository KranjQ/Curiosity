package graph

import (
	"context"
	"curiosity/internal/models"
	"go.uber.org/zap"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

//go:generate go run github.com/99designs/gqlgen

type CommentUseCase interface {
	CreateComment(ctx context.Context, comment models.Comment) error
	GetCommentByID(ctx context.Context, id int) (*models.Comment, error)
	//GetCommentsByPostID(ctx context.Context, postID int, limit int, offset int) ([]*models.Node, error)
	GetCommentsByPostID(ctx context.Context, postID int, limit int, offset int) ([]*models.Comment, error)
	GetRepliesByCommentID(ctx context.Context, commentID int) ([]*models.Comment, error)
}

type PostUseCase interface {
	CreatePost(ctx context.Context, post models.Post) error
	GetPostByID(ctx context.Context, postID int) (*models.Post, error)
	GetPosts(ctx context.Context) ([]*models.Post, error)
}

type UserUseCase interface {
	RegisterUser(ctx context.Context, user models.User) (*models.User, error)
	SignIn(ctx context.Context, user models.User) (int, error)
}

type Resolver struct {
	commentUseCase CommentUseCase
	postUseCase    PostUseCase
	userUseCase    UserUseCase
	logger         *zap.Logger
}

func NewResolver(pu PostUseCase, cu CommentUseCase, uu UserUseCase, logger *zap.Logger) *Resolver {
	return &Resolver{
		commentUseCase: cu,
		postUseCase:    pu,
		userUseCase:    uu,
		logger:         logger,
	}
}
