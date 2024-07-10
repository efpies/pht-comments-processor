package http

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"pht/comments-processor/transport/serializers"
)

type Client struct {
	baseURL    url.URL
	client     *http.Client
	reqSender  RequestSender
	serializer serializers.Serializer
}

type RequestConfigure func(*http.Request) error

func NewClient(optFns ...ClientConfigFunc) (*Client, error) {
	options := DefaultClientConfig()
	for _, o := range optFns {
		if err := o(&options); err != nil {
			return nil, err
		}
	}

	if options.BaseURL.Host == "" {
		return nil, &emptyURLError{}
	}
	if options.RequestSender == nil {
		return nil, &noRequestSenderError{}
	}
	if options.Serializer == nil {
		return nil, &noSerializerError{}
	}

	return &Client{
		baseURL:    options.BaseURL,
		client:     &http.Client{},
		reqSender:  options.RequestSender,
		serializer: options.Serializer,
	}, nil
}

func (c *Client) BaseURL() url.URL {
	return c.baseURL
}

func (c *Client) SendRequest(method string, targetURL url.URL, body any, configure ...RequestConfigure) (response []byte, statusCode int, err error) {
	req, err := c.buildRequest(method, targetURL, body, configure...)
	if err != nil {
		return nil, 0, err
	}

	return c.sendRequest(req)
}

func (c *Client) buildRequest(method string, targetURL url.URL, body any, configure ...RequestConfigure) (*http.Request, error) {
	fullURL := c.baseURL.ResolveReference(&targetURL)

	var bodyReader io.Reader
	if body != nil {
		serializedBody, err := c.serializer.Serialize(body)
		if err != nil {
			return nil, err
		}

		if len(serializedBody) > 0 {
			bodyReader = bytes.NewReader(serializedBody)
		}
	}

	req, err := http.NewRequest(method, fullURL.String(), bodyReader)
	if err != nil {
		return nil, err
	}

	if bodyReader != nil {
		req.Header.Add("Content-Type", c.serializer.ContentType())
	}

	for _, c := range configure {
		if err = c(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

func (c *Client) sendRequest(req *http.Request) (response []byte, statusCode int, err error) {
	return c.reqSender.Send(c.client, req)
}
