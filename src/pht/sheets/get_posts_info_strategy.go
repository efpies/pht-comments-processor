package sheets

import (
	"github.com/samber/lo"
	"pht/comments-processor/pht/model/dto"
	"pht/comments-processor/pht/sheets/adapters"
	"pht/comments-processor/utils"
)

type GetPostsInfoStrategy struct {
	sheetsDataProvider *DataProvider
}

func NewGetPostsInfoStrategy(
	sheetsDataProvider *DataProvider,
) *GetPostsInfoStrategy {
	return &GetPostsInfoStrategy{
		sheetsDataProvider: sheetsDataProvider,
	}
}

func (s *GetPostsInfoStrategy) Fetch(adapter PostAdapter, flowID string, sheet string) (map[string]dto.TablePosts, error) {
	rows, err := s.sheetsDataProvider.GetSheetData(flowID, sheet)
	if err != nil {
		return nil, err
	}

	result := make(map[string]dto.TablePosts)
	if !adapter.IsMultiTable() {
		posts, lastTime, err := s.parseSubTable(adapter, rows)
		if err != nil {
			return nil, err
		}

		result[sheet] = dto.TablePosts{Posts: posts, LastCheckTime: lastTime}
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

			result[header] = dto.TablePosts{Posts: parsed, LastCheckTime: lastTime}

			rows = rows[len(parsed)+1:]
		}
	}

	return result, nil
}

func (s *GetPostsInfoStrategy) parseSubTable(adapter PostAdapter, rows [][]string) (posts []dto.TablePost, lastTime *string, err error) {
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

	lastCheckedPost := lo.MaxBy(postInfos, func(candidate adapters.TablePostInfo, curMax adapters.TablePostInfo) bool {
		if curMax.LastCheckTimeIdx() == nil {
			return true
		}

		if candidate.LastCheckTimeIdx() == nil {
			return false
		}

		return *candidate.LastCheckTimeIdx() > *curMax.LastCheckTimeIdx()
	})

	var lastCheckTime *string
	if lastCheckTimeIdx := lastCheckedPost.LastCheckTimeIdx(); lastCheckTimeIdx != nil {
		lastCheckTime = &timeTable[*lastCheckTimeIdx]
	}

	tablePosts := lo.Map(postInfos, func(info adapters.TablePostInfo, _ int) dto.TablePost {
		id, err1 := adapter.GetPostID(info)
		if err1 != nil {
			err = err1
			return dto.TablePost{}
		}

		return dto.TablePost{
			ID:            id,
			Title:         info.Title,
			CommentsCount: info.LastCommentsCount(),
		}
	})
	if err != nil {
		return nil, nil, err
	}
	return tablePosts, lastCheckTime, nil
}
