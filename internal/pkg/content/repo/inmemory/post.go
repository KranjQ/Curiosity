package inmemory

import (
	"context"
	"curiosity/internal/models"
	"sync"
)

type PostRepository struct {
	posts   map[int]models.Post
	mu      *sync.RWMutex
	counter int
}

func NewPostRepository() *PostRepository {
	return &PostRepository{
		posts:   make(map[int]models.Post),
		mu:      &sync.RWMutex{},
		counter: 0,
	}
}

func (repo *PostRepository) CreatePost(ctx context.Context, post models.Post) error {
	_ = ctx
	repo.mu.Lock()
	repo.counter += 1
	post.ID = repo.counter
	repo.posts[post.ID] = post
	repo.mu.Unlock()
	return nil
}

func (repo *PostRepository) GetPosts(ctx context.Context) ([]models.Post, error) {
	_ = ctx
	var postsArray []models.Post
	repo.mu.RLock()
	for _, value := range repo.posts {
		postsArray = append(postsArray, value)
	}
	repo.mu.RUnlock()
	return postsArray, nil
}

func (repo *PostRepository) GetPostByID(ctx context.Context, postID int) (models.Post, error) {
	_ = ctx
	var post models.Post
	repo.mu.RLock()
	post = repo.posts[postID]
	repo.mu.RUnlock()
	return post, nil
}
