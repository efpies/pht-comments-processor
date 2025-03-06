package dto

type CheckPostResult struct {
	Title            string `json:"title"`
	OldCommentsCount int    `json:"old_comments_count"`
	NewCommentsCount any    `json:"new_comments_count"`
	URL              string `json:"url"`
}

type NotifierData struct {
	Posts         []any   `json:"posts"`
	LastCheckTime *string `json:"last_check_time"`
}
