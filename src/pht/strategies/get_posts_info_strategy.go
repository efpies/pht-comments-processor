package strategies

import (
	"github.com/samber/lo"
	"pht/comments-processor/pht/adapters"
	"pht/comments-processor/pht/services"
	"pht/comments-processor/utils"
)

type TablePostsDto struct {
	Posts         []adapters.TablePost `json:"posts"`
	LastCheckTime *string              `json:"last_check_time"`
}

type GetPostsInfoStrategy struct {
	sheetsDataProvider *services.SheetsDataProvider
}

func NewGetPostsInfoStrategy(
	sheetsDataProvider *services.SheetsDataProvider,
) *GetPostsInfoStrategy {
	return &GetPostsInfoStrategy{
		sheetsDataProvider: sheetsDataProvider,
	}
}

func (s *GetPostsInfoStrategy) Fetch(adapter adapters.PostAdapter, flowId string, sheet string) (map[string]TablePostsDto, error) {
	rows, err := s.sheetsDataProvider.GetSheetData(flowId, sheet)
	if err != nil {
		return nil, err
	}

	result := make(map[string]TablePostsDto)
	if !adapter.IsMultiTable() {
		posts, lastTime, err := s.parseSubTable(adapter, rows)
		if err != nil {
			return nil, err
		}

		result[sheet] = TablePostsDto{Posts: posts, LastCheckTime: lastTime}
	} else {
		panic("Not implemented")
	}

	return result, nil
}

func (s *GetPostsInfoStrategy) parseSubTable(adapter adapters.PostAdapter, rows [][]string) (posts []adapters.TablePost, lastTime *string, err error) {
	postEntries := utils.DropWhile(rows, func(row []string) bool {
		return !adapter.IsPost(row)
	})

	postEntries = utils.TakeWhile(postEntries, func(row []string) bool {
		return adapter.IsPost(row)
	})

	if len(postEntries) == 0 {
		return nil, nil, nil
	}

	_, postsEntriesStartIdx, _ := lo.FindIndexOf(rows, func(row []string) bool {
		return adapter.IsPost(row)
	})

	timeTableRow := rows[postsEntriesStartIdx-1]
	timeTable := adapter.GetTimeTable(timeTableRow)

	postInfos := lo.Map(postEntries, func(row []string, _ int) adapters.TablePostInfo {
		info, err1 := adapter.ToTablePostInfo(row)
		if err1 != nil {
			err = err1
		}

		return info
	})

	if err != nil {
		return nil, nil, err
	}

	lastCheckTime := timeTable[lo.Max(lo.Map(postInfos, func(info adapters.TablePostInfo, _ int) int {
		return info.LastCommentCheckTimeIdx()
	}))]

	tablePosts := lo.Map(postInfos, func(info adapters.TablePostInfo, _ int) adapters.TablePost {
		id, err1 := adapter.GetPostId(info)
		if err1 != nil {
			err = err1
			return adapters.TablePost{}
		}

		return adapters.NewTablePost(info.Title, info.LastCommentsCount(), id)
	})
	if err != nil {
		return nil, nil, err
	}
	return tablePosts, &lastCheckTime, nil
}
