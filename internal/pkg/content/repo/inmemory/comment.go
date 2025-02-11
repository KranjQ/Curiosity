package inmemory

import (
	"context"
	"curiosity/internal/models"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
)

type CommentRepository struct {
	comments map[int]models.Comment
	mu       sync.RWMutex
	counter  int
}

func NewCommentRepository() *CommentRepository {
	return &CommentRepository{
		comments: make(map[int]models.Comment),
		mu:       sync.RWMutex{},
		counter:  0,
	}
}

func (repo *CommentRepository) CreateComment(ctx context.Context, comment models.Comment) error {
	_ = ctx
	repo.mu.Lock()
	repo.counter += 1
	comment.ID = repo.counter
	repo.comments[comment.ID] = comment
	strSlice := strings.Split(comment.Path, ",")
	var intSlice []int
	for _, s := range strSlice {
		num, err := strconv.Atoi(s)
		if err != nil {
			return fmt.Errorf("bad atoi error: %w", err)
		}
		intSlice = append(intSlice, num)
	}
	for _, val := range intSlice {
		temp := repo.comments[val]
		temp.Replies += 1
		repo.comments[val] = temp
	}
	repo.mu.Unlock()
	return nil
}

func (repo *CommentRepository) GetCommentByID(ctx context.Context, commentID int) (models.Comment, error) {
	_ = ctx
	repo.mu.RLock()
	comment, exists := repo.comments[commentID]
	if !exists {
		return models.Comment{}, errors.New("comment doesn't exist")
	}
	repo.mu.RUnlock()
	return comment, nil
}

func (repo *CommentRepository) GetCommentsByPostID(ctx context.Context, postID int, limit int, offset int) ([]models.Comment, error) {
	repo.mu.RLock()
	var comments []models.Comment
	limitCounter := 0
	offsetCounter := 0

	for _, value := range repo.comments {
		if limitCounter >= limit {
			break
		}
		for offsetCounter <= offset {
			offsetCounter += 1
			continue
		}
		if value.Parent == -1 && value.Post == postID {
			comments = append(comments, value)
			limitCounter += 1
		}
	}
	repo.mu.RUnlock()
	return comments, nil
}

func (repo *CommentRepository) GetRepliesByCommentID(ctx context.Context, commentID int) ([]models.Comment, error) {
	_ = ctx
	var replies []models.Comment
	repo.mu.RLock()
	for _, value := range repo.comments {
		if value.Parent == commentID {
			replies = append(replies, value)
		}
	}
	repo.mu.RUnlock()
	return replies, nil
}
