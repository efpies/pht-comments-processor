package adapters

import (
	"github.com/samber/lo"
	"net/url"
	"pht/comments-processor/utils"
	"time"
)

type TablePostInfo struct {
	Date                   time.Time
	Title                  string
	URL                    *url.URL
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
		URL:                    url,
		YesterdayCommentsCount: yesterdayCommentsCount,
		CommentsAtTime:         commentsAtTime,
	}
}

func (e *TablePostInfo) LastCheckTimeIdx() *int {
	if len(e.CommentsAtTime) == 0 {
		return nil
	}

	return utils.Ptr(len(e.CommentsAtTime) - 1)
}

func (e *TablePostInfo) LastCommentsCount() int {
	return lo.LastOr(e.CommentsAtTime, e.YesterdayCommentsCount)
}
