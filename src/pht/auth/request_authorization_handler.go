package auth

import (
	"errors"
	"fmt"
	"net/http"
	httpLocal "pht/comments-processor/transport/http"
)

type RequestAuthorizationHandler struct {
	sender            httpLocal.RequestSender
	authTokenProvider AccessTokenProvider
	tokenRefresher    TokensRefresher
}

func NewRequestAuthorizationHandler(sender httpLocal.RequestSender, authTokenProvider AccessTokenProvider, tokenRefresher TokensRefresher) (*RequestAuthorizationHandler, error) {
	if sender == nil {
		return nil, errors.New("request handler is nil")
	}
	if authTokenProvider == nil {
		return nil, errors.New("tokens provider is nil")
	}

	return &RequestAuthorizationHandler{
		sender:            sender,
		authTokenProvider: authTokenProvider,
		tokenRefresher:    tokenRefresher,
	}, nil
}

func (h *RequestAuthorizationHandler) Send(client *http.Client, req *http.Request) (response []byte, statusCode int, err error) {
	retryCount := 1

	for {
		h.authorize(req)

		response, statusCode, err = h.sender.Send(client, req)
		if err != nil {
			return nil, statusCode, err
		}

		if statusCode == http.StatusUnauthorized && retryCount > 0 {
			err = h.refreshTokens()
			if err != nil {
				return nil, statusCode, err
			}

			retryCount--
			continue
		}

		return
	}
}

func (h *RequestAuthorizationHandler) authorize(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Secret %s", h.authTokenProvider.AccessToken()))
}

func (h *RequestAuthorizationHandler) refreshTokens() error {
	_, err := h.tokenRefresher.RefreshTokens()
	return err
}

func WithAuthorization(tokenProvider AccessTokenProvider, tokenRefresher TokensRefresher) httpLocal.ClientConfigFunc {
	return func(o *httpLocal.ClientConfig) error {
		handler, err := NewRequestAuthorizationHandler(o.RequestSender, tokenProvider, tokenRefresher)
		if err != nil {
			return err
		}

		o.RequestSender = handler
		return nil
	}
}
