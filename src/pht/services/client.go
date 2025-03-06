package services

import (
	"github.com/samber/lo"
	"log"
	"net/http"
	"net/url"
	"pht/comments-processor/pht/auth"
	"pht/comments-processor/pht/config"
	"pht/comments-processor/pht/model"
	"pht/comments-processor/pht/model/dto"
	"pht/comments-processor/pht/transport"
	phtHttp "pht/comments-processor/transport/http"
	"strconv"
)

type ordering string

const (
	createdAtAsc  ordering = "created_at"
	createdAtDesc ordering = "-created_at"
)

type FixedPostsGetter interface {
	GetFixedPosts() ([]dto.Post, error)
}

type PostGetter interface {
	GetPost(postID int) (dto.Post, error)
}

type PostCommentsGetter interface {
	GetPostMostRecentComments(postID int) ([]dto.Comment, error)
	GetLastPostComment(postID int) (*dto.Comment, error)
	GetPostComments(postID int, page int, order ordering) (comments []dto.Comment, hasMore bool, err error)
}

type PagesGetter interface {
	GetPages(from int, toInclusiveOpt *int, list string, sublist *string) (posts []dto.Post, hasMore bool, err error)
}

type WikiGetter interface {
	GetWikis() ([]dto.Wiki, error)
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

func (c *Client) GetFixedPosts() ([]dto.Post, error) {
	log.Println("Loading fixed posts")

	targetURL := url.URL{
		Path: "api/v1/publication/feed/fix-list/",
		RawQuery: url.Values{
			"feed_group": {"physical_transformation"},
		}.Encode(),
	}

	var response []dto.Post
	_, err := c.SendRequest(http.MethodGet, targetURL, nil, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) GetPost(postID int) (dto.Post, error) {
	log.Printf("[%d] Loading post", postID)

	targetURL := (&url.URL{
		Path: "api/v1/publication/retrive",
	}).JoinPath(strconv.Itoa(postID))

	var response dto.Post
	_, err := c.SendRequest(http.MethodGet, *targetURL, nil, &response)
	if err != nil {
		return dto.Post{}, err
	}

	return response, nil
}

func (c *Client) GetPostComments(postID int, page int, order ordering) (comments []dto.Comment, hasMore bool, err error) {
	return PagedLoad(page, page, func(page int) (model.Page[dto.Comment], error) {
		targetURL := (&url.URL{
			Path: "api/v1/parent-comment/list",
			RawQuery: url.Values{
				"page":     {strconv.Itoa(page)},
				"ordering": {string(order)},
			}.Encode(),
		}).JoinPath(strconv.Itoa(postID))

		var response model.Page[dto.Comment]
		_, err := c.SendRequest(http.MethodGet, *targetURL, nil, &response)
		return response, err
	})
}

func (c *Client) GetPostMostRecentComments(postID int) ([]dto.Comment, error) {
	result, _, err := c.GetPostComments(postID, 1, createdAtDesc)
	return result, err
}

func (c *Client) GetLastPostComment(postID int) (*dto.Comment, error) {
	comments, err := c.GetPostMostRecentComments(postID)
	if err != nil {
		return nil, err
	}

	if v, ok := lo.First(comments); ok {
		return &v, nil
	}

	return nil, nil
}

func (c *Client) GetPages(from int, toInclusiveOpt *int, list string, sublist *string) (posts []dto.Post, hasMore bool, err error) {
	slug := list
	if sublist != nil {
		slug += "/" + *sublist
	}

	to := 9999
	if toInclusiveOpt != nil {
		to = *toInclusiveOpt
	}

	log.Printf("[%s] Loading pages from %d to %d", slug, from, to)

	posts, hasMore, err = PagedLoad(from, to, func(page int) (model.Page[dto.Post], error) {
		log.Printf("[%s] Loading page %d", slug, page)

		targetURL := (&url.URL{
			Path: "api/v1/publication",
			RawQuery: url.Values{
				"page":               {strconv.Itoa(page)},
				"feed_group":         {"physical_transformation"},
				"visible_page_count": {"100"},
			}.Encode(),
		}).JoinPath(list).JoinPath("list")

		if sublist != nil {
			targetURL = targetURL.JoinPath(*sublist)
		}

		var subResponse model.Page[dto.Post]
		_, err := c.SendRequest(http.MethodGet, *targetURL, nil, &subResponse)

		return subResponse, err
	})

	if err != nil {
		return nil, false, err
	}

	return posts, hasMore, nil
}

func (c *Client) GetWikis() ([]dto.Wiki, error) {
	log.Println("Loading wikis")

	targetURL := url.URL{
		Path: "api/v1/wiki/page/list/",
	}

	var response []dto.Wiki
	_, err := c.SendRequest(http.MethodGet, targetURL, nil, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
