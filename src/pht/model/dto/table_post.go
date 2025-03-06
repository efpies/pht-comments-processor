package dto

type TablePost struct {
	ID            int    `json:"post_id"`
	Title         string `json:"title"`
	CommentsCount int    `json:"comments_count"`
}
