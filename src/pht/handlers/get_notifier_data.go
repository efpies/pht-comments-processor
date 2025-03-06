package handlers

import (
	"pht/comments-processor/pht/config"
	"pht/comments-processor/pht/model/dto"
	"pht/comments-processor/pht/services"
	"pht/comments-processor/pht/strategies"
)

type getNotifierDataRequest struct {
	PostType PostType `mapstructure:"post_type"`
	FlowID   string   `mapstructure:"flow_id"`
	Sheet    string   `mapstructure:"sheet"`
}

func getNotifierData(sheetsDataProvider *services.SheetsDataProvider, config config.ConfigProvider, checkPostStrategy strategies.CheckPostStrategy) lambdaHandlerInOut[getNotifierDataRequest, map[string]dto.NotifierData] {
	return func(req getNotifierDataRequest) (map[string]dto.NotifierData, error) {
		adapter, err := makePostAdapter(req.PostType, config)
		if err != nil {
			return nil, err
		}

		getPostsInfo := strategies.NewGetPostsInfoStrategy(sheetsDataProvider)
		postsInfo, err := getPostsInfo.Fetch(adapter, req.FlowID, req.Sheet)
		if err != nil {
			return nil, err
		}

		response := make(map[string]dto.NotifierData)
		for sheet, tablePostsDto := range postsInfo {
			var res []any
			for _, post := range tablePostsDto.Posts {
				check, err := checkPostStrategy.CheckPost(post)
				if err == nil {
					res = append(res, check)
				} else {
					res = append(res, dto.Error{Error: err.Error()})
				}
			}

			response[sheet] = dto.NotifierData{
				Posts:         res,
				LastCheckTime: tablePostsDto.LastCheckTime,
			}
		}

		return response, nil
	}
}
