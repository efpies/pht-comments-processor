package adapters

import (
	"pht/comments-processor/pht/config"
	"regexp"
	"strings"
)

type WikiPostAdapter struct {
	postAdapterBase
	contentURI string
}

func NewWikiPostAdapter(config config.ConfigProvider) *WikiPostAdapter {
	const (
		cellIdxTitle                  = 0
		cellIdxLink                   = 1
		cellIdxYesterdayCommentsCount = 2
		cellIdxTimeTableBegin         = 3
	)

	return &WikiPostAdapter{
		postAdapterBase: newPostAdapter(
			withTitleAt(cellIdxTitle),
			withLinkAt(cellIdxLink),
			withYesterdayCommentsCountAt(cellIdxYesterdayCommentsCount),
			withTimeTableStartingAt(cellIdxTimeTableBegin),
			withRegex(regexp.MustCompile(`/publicate/(\d+)`)),
		),
		contentURI: config.ContentURL(),
	}
}

func (a *WikiPostAdapter) IsMultiTable() bool {
	return true
}

func (a *WikiPostAdapter) IsHeader(row []string) bool {
	link := a.getLinkString(row)
	return link != nil && strings.Contains(*link, "/#/listwiki")
}

func (a *WikiPostAdapter) IsPost(row []string) bool {
	link := a.getLinkString(row)
	return link != nil && strings.Contains(*link, a.contentURI) && !a.IsHeader(row)
}
