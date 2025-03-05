package adapters

import (
	"pht/comments-processor/pht/config"
	"regexp"
	"strings"
)

type CommunityPostAdapter struct {
	postAdapterBase
	contentURI string
}

func NewCommunityPostAdapter(config config.ConfigProvider) *CommunityPostAdapter {
	const (
		cellIdxTitle                  = 0
		cellIdxLink                   = 1
		cellIdxYesterdayCommentsCount = 2
		cellIdxTimeTableBegin         = 3
	)

	return &CommunityPostAdapter{
		postAdapterBase: newPostAdapter(
			withTitleAt(cellIdxTitle),
			withLinkAt(cellIdxLink),
			withYesterdayCommentsCountAt(cellIdxYesterdayCommentsCount),
			withTimeTableStartingAt(cellIdxTimeTableBegin),
			withRegex(regexp.MustCompile(`/topic/(\d+)`)),
		),
		contentURI: config.ContentURL(),
	}
}

func (a *CommunityPostAdapter) IsMultiTable() bool {
	return false
}

func (a *CommunityPostAdapter) IsHeader(row []string) bool {
	link := a.getLinkString(row)
	return link != nil && strings.Contains(*link, a.contentURI)
}

func (a *CommunityPostAdapter) IsPost(row []string) bool {
	link := a.getLinkString(row)
	return link != nil && strings.Contains(*link, a.contentURI)
}
