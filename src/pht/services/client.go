package services

import (
	"github.com/samber/lo"
	"log"
	"net/http"
	"net/url"
	"pht/comments-processor/pht/auth"
	"pht/comments-processor/pht/config"
	"pht/comments-processor/pht/model"
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
	GetFixedPosts() ([]model.PostDto, error)
}

type PostGetter interface {
	GetPost(postID int) (model.PostDto, error)
}

type PostCommentsGetter interface {
	GetPostMostRecentComments(postID int) ([]model.CommentDto, error)
	GetLastPostComment(postID int) (*model.CommentDto, error)
	GetPostComments(postID int, page int, order ordering) (comments []model.CommentDto, hasMore bool, err error)
}

type PagesGetter interface {
	GetPages(from int, toInclusiveOpt *int, list string, sublist *string) (posts []model.PostDto, hasMore bool, err error)
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

func (c *Client) GetPost(postID int) (model.PostDto, error) {
	log.Printf("[%d] Loading post", postID)

	targetURL := (&url.URL{
		Path: "api/v1/publication/retrive",
	}).JoinPath(strconv.Itoa(postID))

	var response model.PostDto
	_, err := c.SendRequest(http.MethodGet, *targetURL, nil, &response)
	if err != nil {
		return model.PostDto{}, err
	}

	return response, nil
}

func (c *Client) GetPostComments(postID int, page int, order ordering) (comments []model.CommentDto, hasMore bool, err error) {
	return PagedLoad(page, page, func(page int) (model.Page[model.CommentDto], error) {
		targetURL := (&url.URL{
			Path: "api/v1/parent-comment/list",
			RawQuery: url.Values{
				"page":     {strconv.Itoa(page)},
				"ordering": {string(order)},
			}.Encode(),
		}).JoinPath(strconv.Itoa(postID))

		var response model.Page[model.CommentDto]
		_, err := c.SendRequest(http.MethodGet, *targetURL, nil, &response)
		return response, err
	})
}

func (c *Client) GetPostMostRecentComments(postID int) ([]model.CommentDto, error) {
	result, _, err := c.GetPostComments(postID, 1, createdAtDesc)
	return result, err
}

func (c *Client) GetLastPostComment(postID int) (*model.CommentDto, error) {
	comments, err := c.GetPostMostRecentComments(postID)
	if err != nil {
		return nil, err
	}

	if v, ok := lo.First(comments); ok {
		return &v, nil
	}

	return nil, nil
}

func (c *Client) GetPages(from int, toInclusiveOpt *int, list string, sublist *string) (posts []model.PostDto, hasMore bool, err error) {
	slug := list
	if sublist != nil {
		slug += "/" + *sublist
	}

	to := 9999
	if toInclusiveOpt != nil {
		to = *toInclusiveOpt
	}

	log.Printf("[%s] Loading pages from %d to %d", slug, from, to)

	posts, hasMore, err = PagedLoad(from, to, func(page int) (model.Page[model.PostDto], error) {
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

		var subResponse model.Page[model.PostDto]
		_, err := c.SendRequest(http.MethodGet, *targetURL, nil, &subResponse)

		return subResponse, err
	})

	if err != nil {
		return nil, false, err
	}

	return posts, hasMore, nil
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
