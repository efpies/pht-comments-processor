package handlers

import (
	"pht/comments-processor/pht/model"
	"pht/comments-processor/pht/services"
)

func getFixedPosts(fixedPostsGetter services.FixedPostsGetter) lambdaHandlerOut[[]model.PostDto] {
	return func() ([]model.PostDto, error) {
		return fixedPostsGetter.GetFixedPosts()
	}
}
