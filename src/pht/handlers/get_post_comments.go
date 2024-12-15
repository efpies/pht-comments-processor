package handlers

import (
	"pht/comments-processor/pht/model"
	"pht/comments-processor/pht/services"
)

type getPostCommentsRequest struct {
	PostID int `mapstructure:"post_id"`
}

func getPostComments(postCommentsGetter services.PostCommentsGetter) lambdaHandlerInOut[getPostCommentsRequest, []model.CommentDto] {
	return func(req getPostCommentsRequest) ([]model.CommentDto, error) {
		return postCommentsGetter.GetPostMostRecentComments(req.PostID)
	}
}
