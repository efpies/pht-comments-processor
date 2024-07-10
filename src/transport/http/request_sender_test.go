package http

import (
	"errors"
	"github.com/h2non/gock"
	"net/http"
	"pht/comments-processor/utils"
	"reflect"
	"strings"
	"testing"
)

func Test_DefaultRequestHandler_Send(t *testing.T) {
	testCases := []struct {
		name             string
		wantResponseBody *string
		wantStatusCode   int
		wantErr          error
	}{
		{
			name:             "OK handled correctly",
			wantResponseBody: utils.Ptr("{my body}"),
			wantStatusCode:   http.StatusOK,
		},
		{
			name:             "unauthorized returned as is",
			wantResponseBody: utils.Ptr("{got unauthorized}"),
			wantStatusCode:   http.StatusUnauthorized,
		},
		{
			name:           "no status code on error, error returned as is",
			wantStatusCode: 0,
			wantErr:        errors.New("error happened"),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			defer gock.Off()
			gockReq := gock.New("http://abc.def").
				Get("/qwe").
				Reply(tt.wantStatusCode).
				SetError(tt.wantErr)
			if tt.wantResponseBody != nil {
				gockReq.Body(strings.NewReader(*tt.wantResponseBody))
			}

			h := NewDefaultRequestSender()
			request, _ := http.NewRequest(http.MethodGet, "http://abc.def/qwe", nil)
			gotResponseBody, gotStatusCode, err := h.Send(&http.Client{}, request)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("bad error (got %v, want %v)", err, tt.wantErr)
				t.FailNow()
			}
			if !reflect.DeepEqual(string(gotResponseBody), utils.SafeDeref(tt.wantResponseBody)) {
				t.Errorf("bad responseBody (got %+v, want %+v)", string(gotResponseBody), tt.wantResponseBody)
			}
			if gotStatusCode != tt.wantStatusCode {
				t.Errorf("bad statusCode (got %+v, want %+v)", gotStatusCode, tt.wantStatusCode)
			}
		})
	}
}
