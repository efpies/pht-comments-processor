package sheets

import (
	"pht/comments-processor/pht/config"
	"pht/comments-processor/pht/model/dto"
	"pht/comments-processor/pht/strategies"
)

type NotifierDataGetter struct {
	getPostsInfo      *GetPostsInfoStrategy
	checkPostStrategy strategies.CheckPostStrategy
	config            config.ConfigProvider
}

func NewNotifierDataGetter(
	getPostsInfo *GetPostsInfoStrategy,
	checkPostStrategy strategies.CheckPostStrategy,
	config config.ConfigProvider,
) *NotifierDataGetter {
	return &NotifierDataGetter{
		getPostsInfo:      getPostsInfo,
		checkPostStrategy: checkPostStrategy,
		config:            config,
	}
}

func (g *NotifierDataGetter) GetNotifierData(postType PostType, flowID, sheet string) (map[string]dto.NotifierData, error) {
	adapter, err := MakePostAdapter(postType, g.config)
	if err != nil {
		return nil, err
	}

	postsInfo, err := g.getPostsInfo.Fetch(adapter, flowID, sheet)
	if err != nil {
		return nil, err
	}

	response := make(map[string]dto.NotifierData)
	for section, tablePostsDto := range postsInfo {
		var res []any
		for _, post := range tablePostsDto.Posts {
			check, err := g.checkPostStrategy.CheckPost(post)
			if err == nil {
				res = append(res, check)
			} else {
				res = append(res, dto.Error{Error: err.Error()})
			}
		}

		response[section] = dto.NotifierData{
			Posts:         res,
			LastCheckTime: tablePostsDto.LastCheckTime,
		}
	}

	return response, nil
}

var sheetToType = map[string]PostType{
	"Новые посты":  ContentPostType,
	"Старые посты": ContentPostType,
	"Wiki":         WikiPostType,
	"Комьюнити":    CommunityPostType,
}

var sheetToKey = map[string]string{
	"Новые посты":  "new",
	"Старые посты": "old",
	"Wiki":         "wiki",
	"Комьюнити":    "community",
}

func (g *NotifierDataGetter) GetAllNotifierData(flowID string) (map[string][][]any, error) {
	allData := make(map[string][][]any)
	for sheet, postType := range sheetToType {
		data, err := g.GetNotifierData(postType, flowID, sheet)
		if err != nil {
			return nil, err
		}

		var a [][]any
		for k, v := range data {
			a = append(a, []any{k, v})
		}

		allData[sheetToKey[sheet]] = a
	}

	return allData, nil
}
