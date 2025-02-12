package inmemory

import (
	"context"
	"curiosity/internal/models"
	"errors"
	"slices"
	"sync"
)

//go:generate mockgen -destination=./mocks/mock_postCache.go -package=mocks . PostCache
type PostCache interface {
	Get(key int) (models.Post, bool)
	Set(key int, value models.Post)
	GetAll() []models.Post
}

type PostRepository struct {
	cache   PostCache
	mu      *sync.RWMutex
	counter int
}

func NewPostRepository(pc PostCache) *PostRepository {
	return &PostRepository{
		cache:   pc,
		mu:      &sync.RWMutex{},
		counter: 0,
	}
}

func (repo *PostRepository) CreatePost(ctx context.Context, post models.Post) error {
	_ = ctx
	repo.mu.Lock()
	repo.counter += 1
	post.ID = repo.counter
	repo.mu.Unlock()
	repo.cache.Set(post.ID, post)
	return nil
}

func (repo *PostRepository) GetPosts(ctx context.Context) ([]models.Post, error) {
	_ = ctx
	posts := repo.cache.GetAll()
	slices.SortFunc(posts, func(a, b models.Post) int {
		if a.CreatedAt.After(b.CreatedAt) {
			return -1
		}
		if a.CreatedAt.Before(b.CreatedAt) {
			return 1
		}
		return 0
	})
	return posts, nil
}

func (repo *PostRepository) GetPostByID(ctx context.Context, postID int) (models.Post, error) {
	_ = ctx
	var post models.Post
	post, exists := repo.cache.Get(postID)
	if !exists {
		return models.Post{}, errors.New("Post not found")
	}
	return post, nil
}
