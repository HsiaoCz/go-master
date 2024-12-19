package models

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Author    string    `json:"author"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewPost(title, content, author string, tags []string) *Post {
	now := time.Now()
	return &Post{
		ID:        uuid.New().String(),
		Title:     title,
		Content:   content,
		Author:    author,
		Tags:      tags,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
