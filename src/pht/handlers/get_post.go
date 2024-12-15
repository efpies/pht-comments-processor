package handlers

import (
	"net/http"
	"net/url"
	"pht/comments-processor/pht/auth"
	"pht/comments-processor/pht/config"
	"pht/comments-processor/pht/model"
	"pht/comments-processor/pht/services"
	"pht/comments-processor/pht/transport"
	phtHttp "pht/comments-processor/transport/http"
	"strconv"
)

type getPostRequest struct {
	Id int `mapstructure:"id"`
}

func getPost(config config.ConfigProvider, accessTokenProvider auth.AccessTokenProvider, tokensRefresher auth.TokensRefresher, postCommentsGetter services.PostCommentsGetter) lambdaHandlerInOut[getPostRequest, model.PostDto] {
	return func(req getPostRequest) (model.PostDto, error) {
		client, err := transport.NewHTTPClient(phtHttp.WithBaseURL(config.ContentURL()), auth.WithAuthorization(accessTokenProvider, tokensRefresher))
		if err != nil {
			return model.PostDto{}, err
		}

		targetUrl := (&url.URL{
			Path: "api/v1/publication/retrive",
		}).JoinPath(strconv.Itoa(req.Id))

		var response model.PostDto
		_, err = client.SendRequest(http.MethodGet, *targetUrl, nil, &response)
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
