package http

import (
	"io"
	"net/http"
	"pht/comments-processor/utils"
)

type RequestSender interface {
	Send(client *http.Client, req *http.Request) (response []byte, statusCode int, err error)
}

type DefaultRequestSender struct {
}

func NewDefaultRequestSender() *DefaultRequestSender {
	return &DefaultRequestSender{}
}

func (s *DefaultRequestSender) Send(client *http.Client, req *http.Request) (response []byte, statusCode int, err error) {
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer utils.Close(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}

	return body, resp.StatusCode, nil
}
