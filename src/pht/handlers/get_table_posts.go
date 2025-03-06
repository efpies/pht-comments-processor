package handlers

import (
	"pht/comments-processor/pht/config"
	"pht/comments-processor/pht/model/dto"
	"pht/comments-processor/pht/services"
	"pht/comments-processor/pht/strategies"
)

type getTablePostsRequest struct {
	PostType PostType `mapstructure:"post_type"`
	FlowID   string   `mapstructure:"flow_id"`
	Sheet    string   `mapstructure:"sheet"`
}

func getTablePosts(sheetsDataProvider *services.SheetsDataProvider, config config.ConfigProvider) lambdaHandlerInOut[getTablePostsRequest, map[string]dto.TablePosts] {
	return func(req getTablePostsRequest) (map[string]dto.TablePosts, error) {
		adapter, err := makePostAdapter(req.PostType, config)
		if err != nil {
			return nil, err
		}

		getPostsInfo := strategies.NewGetPostsInfoStrategy(sheetsDataProvider)
		return getPostsInfo.Fetch(adapter, req.FlowID, req.Sheet)
	}
}
