package sheets

import "pht/comments-processor/pht/sheets/adapters"

type PostAdapter interface {
	IsMultiTable() bool
	IsHeader(row []string) bool
	IsPost(row []string) bool
	ToTablePostInfo(row []string) (adapters.TablePostInfo, error)
	GetTimeTable(row []string) []string
	GetPostID(postInfo adapters.TablePostInfo) (int, error)
}
