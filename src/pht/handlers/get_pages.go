package handlers

import (
	"pht/comments-processor/pht/model"
	"pht/comments-processor/pht/services"
)

type getPagesRequest struct {
	From    int     `mapstructure:"from"`
	To      *int    `mapstructure:"to"`
	List    string  `mapstructure:"list"`
	Sublist *string `mapstructure:"sublist"`
}

func getPages(pagesGetter services.PagesGetter, postCommentsGetter services.PostCommentsGetter) lambdaHandlerInOut[getPagesRequest, []model.PostDto] {
	return func(req getPagesRequest) ([]model.PostDto, error) {
		response, _, err := pagesGetter.GetPages(req.From, req.To, req.List, req.Sublist)

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
