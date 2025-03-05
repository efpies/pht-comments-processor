package adapters

import (
	"github.com/samber/lo"
	"net/url"
	"time"
)

type TablePostInfo struct {
	Date                   time.Time
	Title                  string
	Url                    *url.URL
	YesterdayCommentsCount int
	CommentsAtTime         []int
}

func newTablePostInfo(
	date time.Time,
	title string,
	url *url.URL,
	yesterdayCommentsCount int,
	commentsAtTime []int,
) TablePostInfo {
	return TablePostInfo{
		Date:                   date,
		Title:                  title,
		Url:                    url,
		YesterdayCommentsCount: yesterdayCommentsCount,
		CommentsAtTime:         commentsAtTime,
	}
}

func (e *TablePostInfo) LastCommentCheckTimeIdx() int {
	return len(e.CommentsAtTime) - 1
}

func (e *TablePostInfo) LastCommentsCount() int {
	return lo.LastOr(e.CommentsAtTime, e.YesterdayCommentsCount)
}
