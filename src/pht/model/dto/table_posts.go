package dto

type TablePosts struct {
	Posts         []TablePost `json:"posts"`
	LastCheckTime *string     `json:"last_check_time"`
}
