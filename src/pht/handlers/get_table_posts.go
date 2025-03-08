package handlers

import (
	"pht/comments-processor/pht/config"
	"pht/comments-processor/pht/model/dto"
	"pht/comments-processor/pht/sheets"
)

type getTablePostsRequest struct {
	PostType sheets.PostType `mapstructure:"post_type"`
	FlowID   string          `mapstructure:"flow_id"`
	Sheet    string          `mapstructure:"sheet"`
}

func getTablePosts(getPostsInfo *sheets.GetPostsInfoStrategy, config config.ConfigProvider) lambdaHandlerInOut[getTablePostsRequest, map[string]dto.TablePosts] {
	return func(req getTablePostsRequest) (map[string]dto.TablePosts, error) {
		adapter, err := sheets.MakePostAdapter(req.PostType, config)
		if err != nil {
			return nil, err
		}

		return getPostsInfo.Fetch(adapter, req.FlowID, req.Sheet)
	}
}
