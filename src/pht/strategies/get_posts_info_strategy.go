package strategies

import (
	"github.com/samber/lo"
	"pht/comments-processor/pht/adapters"
	"pht/comments-processor/pht/model"
	"pht/comments-processor/pht/services"
	"pht/comments-processor/utils"
)

type TablePostsDto struct {
	Posts         []model.TablePostDto `json:"posts"`
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

func (s *GetPostsInfoStrategy) Fetch(adapter adapters.PostAdapter, flowID string, sheet string) (map[string]TablePostsDto, error) {
	rows, err := s.sheetsDataProvider.GetSheetData(flowID, sheet)
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
		for len(rows) > 0 {
			rows = lo.DropWhile(rows, func(row []string) bool {
				b := !adapter.IsHeader(row)
				return b
			})

			if len(rows) == 0 {
				break
			}

			header := rows[0][0]
			parsed, lastTime, err := s.parseSubTable(adapter, rows)
			if err != nil {
				return nil, err
			}

			result[header] = TablePostsDto{Posts: parsed, LastCheckTime: lastTime}

			rows = rows[len(parsed)+1:]
		}
	}

	return result, nil
}

func (s *GetPostsInfoStrategy) parseSubTable(adapter adapters.PostAdapter, rows [][]string) (posts []model.TablePostDto, lastTime *string, err error) {
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

	lastCheckedPost := lo.MaxBy(postInfos, func(a adapters.TablePostInfo, b adapters.TablePostInfo) bool {
		if a.LastCheckTimeIdx() == nil {
			return true
		}

		if b.LastCheckTimeIdx() == nil {
			return false
		}

		return *a.LastCheckTimeIdx() < *b.LastCheckTimeIdx()
	})

	var lastCheckTime *string
	if lastCheckTimeIdx := lastCheckedPost.LastCheckTimeIdx(); lastCheckTimeIdx != nil {
		lastCheckTime = &timeTable[*lastCheckTimeIdx]
	}

	tablePosts := lo.Map(postInfos, func(info adapters.TablePostInfo, _ int) model.TablePostDto {
		id, err1 := adapter.GetPostID(info)
		if err1 != nil {
			err = err1
			return model.TablePostDto{}
		}

		return model.NewTablePostDto(id, info.Title, info.LastCommentsCount())
	})
	if err != nil {
		return nil, nil, err
	}
	return tablePosts, lastCheckTime, nil
}
