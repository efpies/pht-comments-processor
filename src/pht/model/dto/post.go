package dto

import "time"

type Post struct {
	ID              int       `json:"id"`
	Title           string    `json:"title"`
	Active          *bool     `json:"active"`
	DisableComments bool      `json:"disable_comments"`
	CommentsCount   int       `json:"comments_count"`
	PubStartDate    time.Time `json:"pub_start_date"`
	CreatedAt       time.Time `json:"created_at"`
	LastCommentID   *int      `json:"last_comment_id,omitempty"`
}
