package auth

import (
	"net/http"
	"net/url"
	"pht/comments-processor/pht/config"
	"pht/comments-processor/pht/transport"
	phtHttp "pht/comments-processor/transport/http"
)

type TokensRefresher interface {
	RefreshTokens() (newAccessToken string, err error)
}

type tokensDto struct {
	AccessToken  string `json:"access"`
	RefreshToken string `json:"refresh"`
}

type refreshTokenReq struct {
	RefreshToken string `json:"refresh"`
}

type tokensRefresher struct {
	client               *transport.HTTPClient
	refreshTokenProvider RefreshTokenProvider
	authTokenUpdater     AccessTokenUpdater
	refreshTokenUpdater  RefreshTokenUpdater
}

func NewTokensRefresher(config config.ConfigProvider, refreshTokenProvider RefreshTokenProvider, authTokenUpdater AccessTokenUpdater, refreshTokenUpdater RefreshTokenUpdater) (TokensRefresher, error) {
	client, err := transport.NewHTTPClient(phtHttp.WithBaseURL(config.RefreshTokenURL()))
	if err != nil {
		return nil, err
	}

	return &tokensRefresher{
		client:               client,
		refreshTokenProvider: refreshTokenProvider,
		authTokenUpdater:     authTokenUpdater,
		refreshTokenUpdater:  refreshTokenUpdater,
	}, nil
}

func (c *tokensRefresher) RefreshTokens() (newAccessToken string, err error) {
	newTokens, err := c.getNewTokens()
	if err != nil {
		return "", err
	}

	newAccessToken = newTokens.AccessToken
	err = c.authTokenUpdater.UpdateAccessToken(newAccessToken)
	if err != nil {
		return "", err
	}

	err = c.refreshTokenUpdater.UpdateRefreshToken(newTokens.RefreshToken)
	if err != nil {
		return "", err
	}

	return newAccessToken, nil
}

func (c *tokensRefresher) getNewTokens() (tokensDto, error) {
	req := refreshTokenReq{
		RefreshToken: c.refreshTokenProvider.RefreshToken(),
	}

	newTokens := tokensDto{}

	_, err := c.client.SendRequest(http.MethodPost, url.URL{}, req, &newTokens)
	return newTokens, err
}
