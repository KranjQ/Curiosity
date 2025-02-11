package psql

import (
	"context"
	"curiosity/internal/models"
	"database/sql"
	"fmt"
)

type PostRepository struct {
	DB *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{DB: db}
}

func (repo *PostRepository) CreatePost(ctx context.Context, post models.Post) error {
	query := "INSERT INTO posts (title, content, author, is_commentable) VALUES ($1, $2, $3, $4)"

	_, err := repo.DB.Exec(query, post.Title, post.Content, post.Author, post.IsCommentable)
	if err != nil {
		return fmt.Errorf("bad create post: %w", err)
	}

	return nil
}

func (repo *PostRepository) GetPosts(ctx context.Context) ([]models.Post, error) {
	query := "SELECT id, title, content, author, is_commentable, created_at FROM posts"

	rows, err := repo.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("bad get posts: %w", err)
	}
	defer rows.Close()

	var posts []models.Post

	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content,
			&post.Author, &post.IsCommentable, &post.CreatedAt); err != nil {
			return nil, fmt.Errorf("bad scan in get posts: %w", err)
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (repo *PostRepository) GetPostByID(ctx context.Context, postID int) (models.Post, error) {
	var post models.Post
	query := "SELECT id, title, content, author, is_commentable, created_at  FROM posts WHERE id = $1"
	row := repo.DB.QueryRowContext(ctx, query, postID)
	err := row.Scan(&post.ID, &post.Title, &post.Content, &post.Author, &post.IsCommentable, &post.CreatedAt)
	if err != nil {
		return models.Post{}, fmt.Errorf("bad GetPostByID: %w", err)
	}
	return post, nil
}
