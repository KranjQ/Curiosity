package repo

import (
	"context"
	"curiosity/internal/models"
	"database/sql"
	"fmt"
	"log"
)

type CommentRepository struct {
	DB *sql.DB
}

func NewCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{DB: db}
}

func (repo *CommentRepository) CreateComment(ctx context.Context, comment models.Comment) error {
	tx, err := repo.DB.Begin()
	if err != nil {
		return fmt.Errorf("bad transaction begin error: %w", err)
	}
	defer func() {
		err := tx.Rollback()
		if err != nil && err != sql.ErrTxDone {
			// log
		}
	}()

	insertQuery := "INSERT INTO comments (author, message, post, parent, depth, path) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"

	row := tx.QueryRowContext(ctx, insertQuery, comment.Author, comment.Message,
		comment.Post, comment.Parent, comment.Depth, comment.Path)
	var commentID int
	if err = row.Scan(&commentID); err != nil {
		return fmt.Errorf("bad scan id of new comment error: %w", err)
	}
	log.Print(comment.Path)
	if comment.Path != "" {
		updateRepliesQuery := "UPDATE comments SET replies = replies + 1 WHERE id = ANY(string_to_array($1, ',')::int[])"
		_, err = tx.ExecContext(ctx, updateRepliesQuery, comment.Path)
		if err != nil {
			return fmt.Errorf("bad update path error: %w", err)
		}
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error while tx commit: %w", err)
	}
	return nil
}

func (repo *CommentRepository) GetCommentByID(ctx context.Context, id int) (*models.Comment, error) {
	var comment models.Comment
	query := "SELECT id, author, message, post, parent, depth, path FROM comments WHERE id = $1"
	row := repo.DB.QueryRowContext(ctx, query, id)
	if err := row.Scan(&comment.ID, &comment.Author, &comment.Message, &comment.Post, &comment.Parent,
		&comment.Depth, &comment.Path); err != nil {
		return nil, fmt.Errorf("GetCommentByID error: %w", err)
	}
	return &comment, nil
}

func (repo *CommentRepository) GetCommentsByPostID(ctx context.Context, postID int, limit int, offset int) ([]*models.Comment, error) {
	var comments []*models.Comment
	query := `SELECT id, author, message, post, parent, depth, path, replies, created_at 
FROM comments WHERE post = $1 AND parent = -1 LIMIT $2 OFFSET $3`
	rows, err := repo.DB.QueryContext(ctx, query, postID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("bad GetCommentsByPostID: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var comment models.Comment
		if err := rows.Scan(&comment.ID, &comment.Author, &comment.Message, &comment.Post, &comment.Parent, &comment.Depth,
			&comment.Path, &comment.Replies, &comment.CreatedAt); err != nil {
			return nil, fmt.Errorf("GetCommentsByID bad comment scan: %w", err)
		}
		comments = append(comments, &comment)
	}
	return comments, nil
}

func (repo *CommentRepository) GetRepliesByCommentID(ctx context.Context, commentID int) ([]*models.Comment, error) {
	var comments []*models.Comment
	query := "SELECT id, author, message, post, parent, depth, path, replies, created_at " +
		"FROM comments WHERE parent = $1"
	rows, err := repo.DB.QueryContext(ctx, query, commentID)
	if err != nil {
		return nil, fmt.Errorf("bad GetCommentsByPostID: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var comment models.Comment
		if err := rows.Scan(&comment.ID, &comment.Author, &comment.Message, &comment.Post, &comment.Parent, &comment.Depth,
			&comment.Path, &comment.Replies, &comment.CreatedAt); err != nil {
			return nil, fmt.Errorf("GetCommentsByID bad comment scan: %w", err)
		}
		comments = append(comments, &comment)
	}
	return comments, nil
}
