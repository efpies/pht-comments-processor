package adapters

type TablePost struct {
	Title         string `json:"title"`
	CommentsCount int    `json:"comments_count"`
	PostId        int    `json:"post_id"`
}

func NewTablePost(title string, commentsCount int, postId int) TablePost {
	return TablePost{
		Title:         title,
		CommentsCount: commentsCount,
		PostId:        postId,
	}
}
