package handlers

import (
	"pht/comments-processor/pht/model"
	"pht/comments-processor/pht/services"
)

type getPostRequest struct {
	ID int `mapstructure:"id"`
}

func getPost(postGetter services.PostGetter, postCommentsGetter services.PostCommentsGetter) lambdaHandlerInOut[getPostRequest, model.PostDto] {
	return func(req getPostRequest) (model.PostDto, error) {
		response, err := postGetter.GetPost(req.ID)
		if err != nil {
			return model.PostDto{}, err
		}

		pf := services.NewPostFiller(postCommentsGetter)
		if err := pf.FillPost(&response); err != nil {
			return model.PostDto{}, err
		}

		return response, nil
	}
}
