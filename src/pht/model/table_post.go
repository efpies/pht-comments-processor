package model

type TablePostDto struct {
	ID            int    `json:"post_id"`
	Title         string `json:"title"`
	CommentsCount int    `json:"comments_count"`
}

func NewTablePostDto(id int, title string, commentsCount int) TablePostDto {
	return TablePostDto{
		ID:            id,
		Title:         title,
		CommentsCount: commentsCount,
	}
}
