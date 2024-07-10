package http

import (
	"github.com/h2non/gock"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestClientRequiresBaseURL(t *testing.T) {
	c, err := NewClient()
	if err != nil {
		t.Errorf("error creating a client: %s", err)
	}

	defURL := url.URL{}
	if c.BaseURL() != defURL {
		t.Errorf("base url is not nil: %s", c.baseURL)
	}

	_, _, err = c.SendRequest(http.MethodGet, "/", nil)
	if err == nil {
		t.Errorf("expected an error, but it went ok")
	}
}

func TestClientConfigureRequest(t *testing.T) {
	header := "X_MY_HEADER_FROM_TESTS"
	headerValue := "X_TEST_VALUE"
	expectedBody := "{my body}"

	cases := []struct {
		name         string
		applyHeaders bool
	}{
		{
			name:         "client configures request if requested",
			applyHeaders: true,
		},
		{
			name:         "client doesn't configure request if not requested",
			applyHeaders: false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			defer gock.Off()

			gock.New("https://abc.def").
				MatchHeader(header, headerValue).
				Get("/qwe").
				Reply(200).
				Body(strings.NewReader(expectedBody))

			c, err := NewClient(WithBaseURL("https://abc.def"))
			if err != nil {
				t.Errorf("error creating a client: %s", err)
			}

			body, sc, err := c.SendRequest(http.MethodGet, "/qwe", nil, func(request *http.Request) {
				if tt.applyHeaders {
					request.Header.Add(header, headerValue)
				}
			})

			if tt.applyHeaders && (body == nil || string(body) == expectedBody || sc != 200 || err != nil) ||
				!tt.applyHeaders && (body != nil || sc == 200 || err == nil) {
				t.Errorf("headers not matched (body: %+v, expected: %+v) (status: %d, expected: %d) (err: %s)", string(body), expectedBody, sc, 200, err)
			}
		})
	}
}
