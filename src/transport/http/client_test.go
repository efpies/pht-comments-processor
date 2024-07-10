package http

import (
	"errors"
	"github.com/h2non/gock"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func Test_Configuration_Validated(t *testing.T) {
	t.Parallel()
	configErr := errors.New("configuration error")

	testCases := []struct {
		name      string
		configure ClientConfigFunc
		wantErr   error
	}{
		{
			name:      "correct configuration",
			configure: WithBaseURL("http://localhost"),
		},
		{
			name:      "error during configuration",
			configure: func(c *ClientConfig) error { return configErr },
			wantErr:   configErr,
		},
		{
			name:      "empty base url",
			configure: WithBaseURL(""),
			wantErr:   &emptyURLError{},
		},
		{
			name: "no request sender",
			configure: func(c *ClientConfig) error {
				c.BaseURL = url.URL{Host: "localhost"}
				c.RequestSender = nil
				return nil
			},
			wantErr: &noRequestSenderError{},
		},
		{
			name: "no serializer",
			configure: func(c *ClientConfig) error {
				c.BaseURL = url.URL{Host: "localhost"}
				c.Serializer = nil
				return nil
			},
			wantErr: &noSerializerError{},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewClient(tt.configure)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("expected `%s`, but received `%s`", tt.wantErr, err)
			}
		})
	}
}

func Test_Request_Configured(t *testing.T) {
	t.Parallel()

	header := "X_MY_HEADER_FROM_TESTS"
	headerValue := "X_TEST_VALUE"

	wantResponse := "{my body}"
	wantStatusCode := http.StatusOK

	testCases := []struct {
		name         string
		applyHeaders bool
	}{
		{
			name:         "request configured if requested",
			applyHeaders: true,
		},
		{
			name:         "request is not configured if not requested",
			applyHeaders: false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			defer gock.Off()

			gock.New("http://localhost").
				MatchHeader(header, headerValue).
				Get("qwe").
				Reply(wantStatusCode).
				Body(strings.NewReader(wantResponse))

			c, err := NewClient(WithBaseURL("http://localhost"))
			if err != nil {
				t.Errorf("error creating a client: %s", err)
				t.FailNow()
			}

			gotResponse, gotStatusCode, err := c.SendRequest(http.MethodGet, url.URL{Path: "qwe"}, nil, func(request *http.Request) error {
				if tt.applyHeaders {
					request.Header.Add(header, headerValue)
				}
				return nil
			})

			if tt.applyHeaders && (gotResponse == nil || string(gotResponse) != wantResponse || gotStatusCode != wantStatusCode || err != nil) ||
				!tt.applyHeaders && (gotResponse != nil || gotStatusCode == wantStatusCode || err == nil) {
				t.Errorf("headers not matched. response: got: %+v, want: %+v), statusCode (got: %d, want: %d) (err: %s)", string(gotResponse), wantResponse, gotStatusCode, 200, err)
			}
		})
	}
}
