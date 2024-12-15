package handlers

import (
	"pht/comments-processor/pht/model"
	"pht/comments-processor/pht/services"
)

func getFixedPosts(fixedPostsGetter services.FixedPostsGetter, postCommentsGetter services.PostCommentsGetter) lambdaHandlerOut[[]model.PostDto] {
	return func() ([]model.PostDto, error) {
		response, err := fixedPostsGetter.GetFixedPosts()
		if err != nil {
			return nil, err
		}

		pf := services.NewPostFiller(postCommentsGetter)
		for i := range response {
			post := &response[i]
			if err := pf.FillPost(post); err != nil {
				return nil, err
			}
		}

		return response, nil
	}
}
