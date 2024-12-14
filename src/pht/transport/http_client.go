package transport

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"pht/comments-processor/pht/model"
	httpTransport "pht/comments-processor/transport/http"
	"pht/comments-processor/utils"
)

type HTTPClient struct {
	httpClient *httpTransport.Client
	hostURL    url.URL
}

func NewHTTPClient(optFns ...httpTransport.ClientConfigFunc) (*HTTPClient, error) {
	httpClient, err := httpTransport.NewClient(optFns...)
	if err != nil {
		return nil, err
	}

	baseURL := httpClient.BaseURL()
	hostURL, err := utils.HostURL(&baseURL)
	if err != nil {
		return nil, err
	}

	return &HTTPClient{
		httpClient: httpClient,
		hostURL:    *hostURL,
	}, nil
}

func (c *HTTPClient) SendRequest(method string, targetUrl url.URL, body any, response any) (statusCode int, err error) {
	respBody, statusCode, err := c.httpClient.SendRequest(method, targetUrl, body, c.withHeaders)
	if err != nil {
		return
	}

	switch statusCode {
	case http.StatusOK:
		err = json.Unmarshal(respBody, response)
	default:
		e := &model.ErrorResponse{}
		if err = json.Unmarshal(respBody, e); err == nil {
			err = e.Error()
		}
	}

	return
}

func (c *HTTPClient) withHeaders(req *http.Request) error {
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Accept-Language", "ru-RU,ru")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("DNT", "1")
	req.Header.Add("Referer", c.hostURL.String())
	req.Header.Add("User-Agent", "comments-processor/"+os.Getenv("APP_VERSION"))
	return nil
}
