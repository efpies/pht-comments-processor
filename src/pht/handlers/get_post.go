package handlers

import (
	"pht/comments-processor/pht/model/dto"
	"pht/comments-processor/pht/services"
)

type getPostRequest struct {
	ID int `mapstructure:"id"`
}

func getPost(postGetter services.PostGetter) lambdaHandlerInOut[getPostRequest, dto.Post] {
	return func(req getPostRequest) (dto.Post, error) {
		return postGetter.GetPost(req.ID)
	}
}
