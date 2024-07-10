package http

import (
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	baseURL   *url.URL
	client    *http.Client
	reqSender RequestSender
}

type RequestConfigure func(*http.Request) error

func NewClient(optFns ...func(*ClientOptions) error) (*Client, error) {
	options := DefaultClientOptions()
	for _, o := range optFns {
		if err := o(options); err != nil {
			return nil, err
		}
	}

	return &Client{
		baseURL:   &options.BaseURL,
		client:    &http.Client{},
		reqSender: options.RequestSender,
	}, nil
}

func (c *Client) BaseURL() url.URL {
	return *c.baseURL
}

func (c *Client) SendRequest(method string, uri string, body io.Reader, configure ...func(*http.Request)) (resp []byte, status int, err error) {
	req, err := c.buildRequest(method, uri, body, configure...)
	if err != nil {
		return nil, 0, err
	}

	return c.sendRequest(req)
}

func (c *Client) buildRequest(method string, uri string, body io.Reader, configure ...func(*http.Request)) (request *http.Request, err error) {
	fullURL := c.baseURL.JoinPath(uri)

	req, err := http.NewRequest(method, fullURL.String(), body)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	for _, c := range configure {
		c(req)
	}

	return req, nil
}

func (c *Client) sendRequest(req *http.Request) (response []byte, status int, err error) {
	return c.reqSender.Send(c.client, req)
}
