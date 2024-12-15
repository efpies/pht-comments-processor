package services

import (
	"log"
	"net/http"
	"net/url"
	"pht/comments-processor/pht/auth"
	"pht/comments-processor/pht/config"
	"pht/comments-processor/pht/model"
	"pht/comments-processor/pht/transport"
	phtHttp "pht/comments-processor/transport/http"
)

type FixedPostsGetter interface {
	GetFixedPosts() ([]model.PostDto, error)
}

type WikiGetter interface {
	GetWikis() ([]model.WikiDto, error)
}

type Client struct {
	*transport.HTTPClient
}

func NewClient(config config.ConfigProvider, accessTokenProvider auth.AccessTokenProvider, tokensRefresher auth.TokensRefresher) (*Client, error) {
	httpClient, err := transport.NewHTTPClient(phtHttp.WithBaseURL(config.ContentURL()), auth.WithAuthorization(accessTokenProvider, tokensRefresher))
	if err != nil {
		return nil, err
	}

	return &Client{
		HTTPClient: httpClient,
	}, nil
}

func (c *Client) GetFixedPosts() ([]model.PostDto, error) {
	log.Println("Loading fixed posts")

	targetURL := url.URL{
		Path: "api/v1/publication/feed/fix-list/",
		RawQuery: url.Values{
			"feed_group": {"physical_transformation"},
		}.Encode(),
	}

	var response []model.PostDto
	_, err := c.SendRequest(http.MethodGet, targetURL, nil, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) GetWikis() ([]model.WikiDto, error) {
	log.Println("Loading wikis")

	targetURL := url.URL{
		Path: "api/v1/wiki/page/list/",
	}

	var response []model.WikiDto
	_, err := c.SendRequest(http.MethodGet, targetURL, nil, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
