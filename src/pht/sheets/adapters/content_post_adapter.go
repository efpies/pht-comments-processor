package adapters

import (
	"pht/comments-processor/pht/config"
	"regexp"
	"strings"
)

type ContentPostAdapter struct {
	postAdapterBase
	contentURI string
}

func NewContentPostAdapter(config config.ConfigProvider) *ContentPostAdapter {
	const (
		cellIdxDate                   = 0
		cellIdxTitle                  = 1
		cellIdxLink                   = 2
		cellIdxYesterdayCommentsCount = 3
		cellIdxTimeTableBegin         = 4
	)

	return &ContentPostAdapter{
		postAdapterBase: newPostAdapter(
			withDateAt(cellIdxDate),
			withTitleAt(cellIdxTitle),
			withLinkAt(cellIdxLink),
			withYesterdayCommentsCountAt(cellIdxYesterdayCommentsCount),
			withTimeTableStartingAt(cellIdxTimeTableBegin),
			withRegex(regexp.MustCompile(`/publicate/(\d+)`)),
		),
		contentURI: config.ContentURL(),
	}
}

func (a *ContentPostAdapter) IsMultiTable() bool {
	return false
}

func (a *ContentPostAdapter) IsHeader(row []string) bool {
	link := a.getLinkString(row)
	return link != nil && strings.Contains(*link, a.contentURI)
}

func (a *ContentPostAdapter) IsPost(row []string) bool {
	link := a.getLinkString(row)
	return link != nil && strings.Contains(*link, a.contentURI)
}
