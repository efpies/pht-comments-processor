package model

type Page[T any] struct {
	Pagination struct {
		PageNumber     int  `json:"page"`
		NextPageNumber *int `json:"next_page_number"`
		HasNextPage    bool `json:"has_next"`
		CountPages     int  `json:"count_pages"`
		TotalCount     int  `json:"total_obj_count"`
	} `json:"pagination"`
	Items []T `json:"results"`
}
