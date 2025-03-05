package adapters

type PostAdapter interface {
	IsMultiTable() bool
	IsHeader(row []string) bool
	IsPost(row []string) bool
	ToTablePostInfo(row []string) (TablePostInfo, error)
	GetTimeTable(row []string) []string
	GetPostId(postInfo TablePostInfo) (int, error)
}
