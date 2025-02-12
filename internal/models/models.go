package models

import "time"

type Post struct {
	ID            int       `json:"id"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	Author        int       `json:"author"`
	IsCommentable bool      `json:"is_commentable"`
	CreatedAt     time.Time `json:"created_at"`
}

type Comment struct {
	ID        int       `json:"id"`
	Author    int       `json:"author"`
	Message   string    `json:"message"`
	Post      int       `json:"post"`
	Parent    int       `json:"parent"`
	Depth     int       `json:"depth"`
	Path      string    `json:"path"`
	Replies   int       `json:"replies"`
	CreatedAt time.Time `json:"created_at"`
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}
