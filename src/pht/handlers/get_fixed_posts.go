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
)

func getFixedPosts(config config.ConfigProvider, accessTokenProvider auth.AccessTokenProvider, tokensRefresher auth.TokensRefresher, postCommentsGetter services.PostCommentsGetter) lambdaHandlerOut[[]model.PostDto] {
	return func() ([]model.PostDto, error) {
		client, err := transport.NewHTTPClient(phtHttp.WithBaseURL(config.ContentURL()), auth.WithAuthorization(accessTokenProvider, tokensRefresher))
		if err != nil {
			return nil, err
		}

		targetUrl := url.URL{
			Path: "api/v1/publication/feed/fix-list/",
			RawQuery: url.Values{
				"feed_group": {"physical_transformation"},
			}.Encode(),
		}

		var response []model.PostDto
		_, err = client.SendRequest(http.MethodGet, targetUrl, nil, &response)
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
