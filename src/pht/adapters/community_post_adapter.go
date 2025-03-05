package adapters

import (
	"fmt"
	"github.com/samber/lo"
	"net/url"
	"pht/comments-processor/pht/config"
	"pht/comments-processor/utils"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type CommunityPostAdapter struct {
	contentURI string
}

const (
	cellIdxTitle                  = 0
	cellIdxLink                   = 1
	cellIdxYesterdayCommentsCount = 2
	cellIdxTimeTableBegin         = 3
)

func NewCommunityPostAdapter(config config.ConfigProvider) *CommunityPostAdapter {
	return &CommunityPostAdapter{contentURI: config.ContentURL()}
}

func (a *CommunityPostAdapter) IsMultiTable() bool {
	return false
}

func (a *CommunityPostAdapter) IsHeader(row []string) bool {
	return len(row) > cellIdxLink && strings.Contains(row[cellIdxLink], a.contentURI)
}

func (a *CommunityPostAdapter) IsPost(row []string) bool {
	return len(row) > cellIdxLink && strings.Contains(row[cellIdxLink], a.contentURI)
}

func (a *CommunityPostAdapter) ToTablePostInfo(row []string) (TablePostInfo, error) {
	if len(row) <= cellIdxLink {
		return TablePostInfo{}, fmt.Errorf("invalid row: %v", row)
	}

	postUrl, err := url.Parse(row[cellIdxLink])
	if err != nil {
		return TablePostInfo{}, err
	}

	yesterdayCommentsCount, err := a.getYesterdayCommentsCount(row)
	if err != nil {
		return TablePostInfo{}, err
	}

	commentsAtTime := lo.Map(a.GetTimeTable(row), func(cell string, _ int) int {
		value, err1 := utils.AtoiSafe(cell)
		if err1 != nil {
			err = err1
			return -1
		}

		return value
	})

	if err != nil {
		return TablePostInfo{}, err
	}

	return newTablePostInfo(
		time.Time{},
		row[cellIdxTitle],
		postUrl,
		yesterdayCommentsCount,
		commentsAtTime,
	), nil
}

func (a *CommunityPostAdapter) getYesterdayCommentsCount(row []string) (int, error) {
	if len(row) <= cellIdxYesterdayCommentsCount {
		return 0, nil
	}

	return utils.AtoiSafe(row[cellIdxYesterdayCommentsCount])
}

func (a *CommunityPostAdapter) GetTimeTable(row []string) []string {
	if len(row) <= cellIdxTimeTableBegin {
		return nil
	}

	return row[cellIdxTimeTableBegin:]
}

var postUrlRegexp = regexp.MustCompile(`/topic/(\d+)`)

func (a *CommunityPostAdapter) GetPostId(post TablePostInfo) (int, error) {
	matches := postUrlRegexp.FindStringSubmatch(post.Url.Fragment)
	if len(matches) < 2 {
		return -1, fmt.Errorf("invalid link: %s", post.Url.Fragment)
	}

	id, err := strconv.Atoi(matches[1])
	if err != nil {
		return -1, err
	}

	return id, nil
}
