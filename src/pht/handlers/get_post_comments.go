package handlers

import (
	"pht/comments-processor/pht/model/dto"
	"pht/comments-processor/pht/services"
)

type getPostCommentsRequest struct {
	PostID int `mapstructure:"post_id"`
}

func getPostComments(postCommentsGetter services.PostCommentsGetter) lambdaHandlerInOut[getPostCommentsRequest, []dto.Comment] {
	return func(req getPostCommentsRequest) ([]dto.Comment, error) {
		return postCommentsGetter.GetPostMostRecentComments(req.PostID)
	}
}
