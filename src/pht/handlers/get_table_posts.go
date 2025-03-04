package handlers

import (
	"fmt"
	"pht/comments-processor/pht/adapters"
	"pht/comments-processor/pht/config"
	"pht/comments-processor/pht/services"
	"pht/comments-processor/pht/strategies"
)

type PostType string

const (
	CommunityPostType PostType = "community"
)

func makePostAdapter(postType PostType, config config.ConfigProvider) (adapters.PostAdapter, error) {
	switch postType {
	case CommunityPostType:
		return adapters.NewCommunityPostAdapter(config), nil
	default:
		return nil, fmt.Errorf("unsupported post type: %s", postType)
	}
}

type getTablePostsRequest struct {
	PostType PostType `mapstructure:"post_type"`
	FlowId   string   `mapstructure:"flow_id"`
	Sheet    string   `mapstructure:"sheet"`
}

func getTablePosts(sheetsDataProvider *services.SheetsDataProvider, config config.ConfigProvider) lambdaHandlerInOut[getTablePostsRequest, map[string]strategies.TablePostsDto] {
	return func(req getTablePostsRequest) (map[string]strategies.TablePostsDto, error) {
		adapter, err := makePostAdapter(req.PostType, config)
		if err != nil {
			return nil, err
		}

		getPostsInfo := strategies.NewGetPostsInfoStrategy(sheetsDataProvider)
		return getPostsInfo.Fetch(adapter, req.FlowId, req.Sheet)
	}
}
