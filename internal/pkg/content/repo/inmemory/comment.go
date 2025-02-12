package inmemory

import (
	"context"
	"curiosity/internal/models"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
)

//go:generate mockgen -destination=./mocks/mock_commentCache.go -package=mocks . CommentCache
type CommentCache interface {
	Get(key int) (models.Comment, bool)
	Set(key int, value models.Comment)
	GetAll() []models.Comment
}

type CommentRepository struct {
	mu      *sync.RWMutex
	cache   CommentCache
	counter int
}

func NewCommentRepository(cc CommentCache) *CommentRepository {
	return &CommentRepository{
		cache:   cc,
		mu:      &sync.RWMutex{},
		counter: 0,
	}
}

func (repo *CommentRepository) CreateComment(ctx context.Context, comment models.Comment) error {
	_ = ctx
	repo.mu.Lock()
	repo.counter += 1
	comment.ID = repo.counter
	repo.mu.Unlock()
	comment.CreatedAt = time.Now()
	repo.cache.Set(comment.ID, comment)
	var intSlice []int
	if comment.Path != "" {
		strSlice := strings.Split(comment.Path, ",")
		for _, s := range strSlice {
			num, err := strconv.Atoi(s)
			if err != nil {
				return fmt.Errorf("bad atoi error: %w", err)
			}
			intSlice = append(intSlice, num)
		}
	}

	for _, val := range intSlice {
		temp, _ := repo.cache.Get(val)
		temp.Replies += 1
		repo.cache.Set(val, temp)
	}

	return nil
}

func (repo *CommentRepository) GetCommentByID(ctx context.Context, commentID int) (models.Comment, error) {
	_ = ctx
	repo.mu.RLock()
	comment, exists := repo.cache.Get(commentID)
	if !exists {
		return models.Comment{}, errors.New("comment doesn't exist")
	}
	repo.mu.RUnlock()
	return comment, nil
}

func (repo *CommentRepository) GetCommentsByPostID(ctx context.Context, postID int, limit int, offset int) ([]models.Comment, error) {
	var filtered []models.Comment

	comments := repo.cache.GetAll()
	slices.SortFunc(comments, func(a, b models.Comment) int {
		if a.CreatedAt.After(b.CreatedAt) {
			return -1
		}
		if a.CreatedAt.Before(b.CreatedAt) {
			return 1
		}
		return 0
	})

	for _, value := range comments {
		if value.Parent == -1 && value.Post == postID {
			filtered = append(filtered, value)
		}
	}

	if offset > len(filtered) {
		return []models.Comment{}, nil
	}
	end := offset + limit
	if end > len(filtered) {
		end = len(filtered)
	}

	return filtered[offset:end], nil
}

func (repo *CommentRepository) GetRepliesByCommentID(ctx context.Context, commentID int) ([]models.Comment, error) {
	_ = ctx
	var replies []models.Comment

	comments := repo.cache.GetAll()
	slices.SortFunc(comments, func(a, b models.Comment) int {
		if a.CreatedAt.After(b.CreatedAt) {
			return -1
		}
		if a.CreatedAt.Before(b.CreatedAt) {
			return 1
		}
		return 0
	})
	for _, value := range comments {
		if value.Parent == commentID {
			replies = append(replies, value)
		}
	}
	return replies, nil
}
