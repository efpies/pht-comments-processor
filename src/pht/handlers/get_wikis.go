package handlers

import (
	"net/http"
	"net/url"
	"pht/comments-processor/pht/auth"
	"pht/comments-processor/pht/config"
	"pht/comments-processor/pht/transport"
	phtHttp "pht/comments-processor/transport/http"
)

type wikiDto struct {
	Id    int    `json:"pk"`
	Title string `json:"title"`
}

func getWikis(config config.ConfigProvider, accessTokenProvider auth.AccessTokenProvider, tokensRefresher auth.TokensRefresher) lambdaHandlerOut[[]wikiDto] {
	return func() ([]wikiDto, error) {
		client, err := transport.NewHTTPClient(phtHttp.WithBaseURL(config.ContentURL()), auth.WithAuthorization(accessTokenProvider, tokensRefresher))
		if err != nil {
			return nil, err
		}

		targetUrl := url.URL{
			Path: "api/v1/wiki/page/list/",
		}

		var response []wikiDto
		_, err = client.SendRequest(http.MethodGet, targetUrl, nil, &response)
		if err != nil {
			return nil, err
		}

		return response, nil
	}
}
