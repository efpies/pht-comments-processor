package http

import "net/url"

type ClientOptions struct {
	BaseURL       url.URL
	RequestSender RequestSender
}

type ClientOptionsConfig func(*ClientOptions) error

func DefaultClientOptions() *ClientOptions {
	return &ClientOptions{
		RequestSender: NewDefaultRequestSender(),
	}
}

func WithBaseURL(baseURL string) ClientOptionsConfig {
	return func(o *ClientOptions) error {
		parsedURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}

		o.BaseURL = *parsedURL
		return nil
	}
}
