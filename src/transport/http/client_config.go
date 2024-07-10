package http

import (
	"net/url"
	"pht/comments-processor/transport/serializers"
)

type ClientConfig struct {
	BaseURL       url.URL
	RequestSender RequestSender
	Serializer    serializers.Serializer
}

type ClientConfigFunc = func(*ClientConfig) error

func DefaultClientConfig() ClientConfig {
	return ClientConfig{
		RequestSender: NewDefaultRequestSender(),
		Serializer:    serializers.NewJsonSerializer(),
	}
}

func WithBaseURL(baseURL string) ClientConfigFunc {
	return func(o *ClientConfig) error {
		parsedURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}

		o.BaseURL = *parsedURL
		return nil
	}
}

func WithSerializer(serializer serializers.Serializer) ClientConfigFunc {
	return func(o *ClientConfig) error {
		o.Serializer = serializer
		return nil
	}
}
