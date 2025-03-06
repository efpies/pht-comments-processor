package handlers

import (
	"pht/comments-processor/pht/model/dto"
	"pht/comments-processor/pht/services"
)

func getFixedPosts(fixedPostsGetter services.FixedPostsGetter) lambdaHandlerOut[[]dto.Post] {
	return func() ([]dto.Post, error) {
		return fixedPostsGetter.GetFixedPosts()
	}
}
