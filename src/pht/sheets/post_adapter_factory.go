package sheets

import (
	"fmt"
	"pht/comments-processor/pht/config"
	"pht/comments-processor/pht/sheets/adapters"
)

type PostType string

const (
	ContentPostType   PostType = "content"
	WikiPostType      PostType = "wiki"
	CommunityPostType PostType = "community"
)

func MakePostAdapter(postType PostType, config config.ConfigProvider) (PostAdapter, error) {
	switch postType {
	case ContentPostType:
		return adapters.NewContentPostAdapter(config), nil
	case WikiPostType:
		return adapters.NewWikiPostAdapter(config), nil
	case CommunityPostType:
		return adapters.NewCommunityPostAdapter(config), nil
	default:
		return nil, fmt.Errorf("unsupported post type: %s", postType)
	}
}
