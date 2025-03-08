package adapters

import (
	"fmt"
	"github.com/samber/lo"
	"net/url"
	"pht/comments-processor/utils"
	"regexp"
	"strconv"
	"time"
)

type adapterConfig struct {
	dateCellIdx                   *int
	titleCellIdx                  *int
	linkCellIdx                   *int
	yesterdayCommentsCountCellIdx *int
	timeTableBeginCellIdx         *int
	postUrlRegexp                 *regexp.Regexp
}

type adapterConfigOption func(*adapterConfig)

type postAdapterBase struct {
	adapterConfig
}

func newPostAdapter(options ...adapterConfigOption) postAdapterBase {
	cfg := adapterConfig{}
	for _, opt := range options {
		opt(&cfg)
	}

	return postAdapterBase{
		adapterConfig: cfg,
	}
}

func (a *postAdapterBase) ToTablePostInfo(row []string) (TablePostInfo, error) {
	link, err := a.getLink(row)
	if link == nil && err == nil {
		err = fmt.Errorf("invalid row: %v", row)
	}
	if err != nil {
		return TablePostInfo{}, err
	}

	yesterdayCommentsCount, err := a.getYesterdayCommentsCount(row)
	if err != nil {
		return TablePostInfo{}, err
	}

	commentsAtTime, err := a.getParsedTimeTable(row)
	if err != nil {
		return TablePostInfo{}, err
	}

	date, err := a.getDate(row)
	if err != nil {
		return TablePostInfo{}, err
	}

	return newTablePostInfo(
		date,
		a.getTitle(row),
		link,
		yesterdayCommentsCount,
		commentsAtTime,
	), nil
}

func (a *postAdapterBase) GetTimeTable(row []string) []string {
	return applyOr(row, a.timeTableBeginCellIdx, func(row []string, idx int) []string {
		return row[idx:]
	}, nil)
}

func (a *postAdapterBase) GetPostID(post TablePostInfo) (int, error) {
	if a.postUrlRegexp == nil {
		return -1, fmt.Errorf("no pattern for post URL")
	}

	matches := a.postUrlRegexp.FindStringSubmatch(post.URL.Fragment)
	if len(matches) < 2 {
		return -1, fmt.Errorf("invalid link: %s", post.URL.Fragment)
	}

	id, err := strconv.Atoi(matches[1])
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (a *postAdapterBase) getDate(row []string) (time.Time, error) {
	return valueOr(row, a.dateCellIdx, func(value string) (time.Time, error) {
		return time.Parse(`02.01.2006`, value)
	}, time.Time{})
}

func (a *postAdapterBase) getTitle(row []string) string {
	result, _ := valueOr(row, a.titleCellIdx, func(value string) (string, error) {
		return value, nil
	}, "")
	return result
}

func (a *postAdapterBase) getLinkString(row []string) *string {
	result, _ := valueOr(row, a.linkCellIdx, func(value string) (*string, error) {
		return &value, nil
	}, nil)
	return result
}

func (a *postAdapterBase) getLink(row []string) (*url.URL, error) {
	link := a.getLinkString(row)
	if link == nil {
		return nil, nil
	}

	return url.Parse(*link)
}

func (a *postAdapterBase) getYesterdayCommentsCount(row []string) (int, error) {
	return valueOr(row, a.yesterdayCommentsCountCellIdx, func(value string) (int, error) {
		return utils.AtoiSafe(value)
	}, 0)
}

func (a *postAdapterBase) getParsedTimeTable(row []string) ([]int, error) {
	var err error
	timeTable := lo.Map(a.GetTimeTable(row), func(cell string, _ int) int {
		value, err1 := utils.AtoiSafe(cell)
		if err1 != nil {
			err = err1
			return -1
		}

		return value
	})
	if err != nil {
		return nil, err
	}

	return timeTable, nil
}

func valueOr[T any](row []string, idx *int, callback func(string) (T, error), defaultValue T) (T, error) {
	if idx == nil || *idx >= len(row) {
		return defaultValue, nil
	}

	return callback(row[*idx])
}

func applyOr[T any](row []string, idx *int, callback func([]string, int) T, defaultValue T) T {
	if idx == nil || *idx >= len(row) {
		return defaultValue
	}

	return callback(row, *idx)
}

func withDateAt(dateCellIdx int) adapterConfigOption {
	return func(cfg *adapterConfig) {
		cfg.dateCellIdx = &dateCellIdx
	}
}

func withTitleAt(titleCellIdx int) adapterConfigOption {
	return func(cfg *adapterConfig) {
		cfg.titleCellIdx = &titleCellIdx
	}
}

func withLinkAt(linkCellIdx int) adapterConfigOption {
	return func(cfg *adapterConfig) {
		cfg.linkCellIdx = &linkCellIdx
	}
}

func withYesterdayCommentsCountAt(yesterdayCommentsCountCellIdx int) adapterConfigOption {
	return func(cfg *adapterConfig) {
		cfg.yesterdayCommentsCountCellIdx = &yesterdayCommentsCountCellIdx
	}
}

func withTimeTableStartingAt(timeTableBeginCellIdx int) adapterConfigOption {
	return func(cfg *adapterConfig) {
		cfg.timeTableBeginCellIdx = &timeTableBeginCellIdx
	}
}

func withRegex(pattern *regexp.Regexp) adapterConfigOption {
	return func(cfg *adapterConfig) {
		cfg.postUrlRegexp = pattern
	}
}
