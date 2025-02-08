package handlers

import (
	"pht/comments-processor/pht/model"
	"pht/comments-processor/pht/services"
)

type getPostRequest struct {
	Id int `mapstructure:"id"`
}

func getPost(postGetter services.PostGetter) lambdaHandlerInOut[getPostRequest, model.PostDto] {
	return func(req getPostRequest) (model.PostDto, error) {
		return postGetter.GetPost(req.Id)
	}
}
