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

type getPagesRequest struct {
	From    int     `mapstructure:"from"`
	To      *int    `mapstructure:"to"`
	List    string  `mapstructure:"list"`
	Sublist *string `mapstructure:"sublist"`
}

func getPages(config config.ConfigProvider, accessTokenProvider auth.AccessTokenProvider, tokensRefresher auth.TokensRefresher, postCommentsGetter services.PostCommentsGetter) lambdaHandlerInOut[getPagesRequest, []model.PostDto] {
	return func(req getPagesRequest) ([]model.PostDto, error) {
		client, err := transport.NewHTTPClient(phtHttp.WithBaseURL(config.ContentURL()), auth.WithAuthorization(accessTokenProvider, tokensRefresher))
		if err != nil {
			return nil, err
		}

		from := req.From
		to := 9999
		if req.To != nil {
			to = *req.To
		}

		response, err := services.PagedLoad(from, to, func(page int) (model.Page[model.PostDto], error) {
			targetUrl := (&url.URL{
				Path: "api/v1/publication",
				RawQuery: url.Values{
					"page":               {strconv.Itoa(page)},
					"feed_group":         {"physical_transformation"},
					"visible_page_count": {strconv.Itoa(100)},
				}.Encode(),
			}).JoinPath(req.List).JoinPath("list")

			if req.Sublist != nil {
				targetUrl = targetUrl.JoinPath(*req.Sublist)
			}

			var subResponse model.Page[model.PostDto]
			_, err = client.SendRequest(http.MethodGet, *targetUrl, nil, &subResponse)

			return subResponse, err
		})

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
