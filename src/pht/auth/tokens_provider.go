package auth

import (
	"fmt"
	"pht/comments-processor/repo"
)

const accessTokenPath = "pht/tokens/access"
const refreshTokenPath = "pht/tokens/refresh"

type AccessTokenProvider interface {
	AccessToken() string
}

type AccessTokenUpdater interface {
	UpdateAccessToken(newToken string) error
}

type RefreshTokenProvider interface {
	RefreshToken() string
}

type RefreshTokenUpdater interface {
	UpdateRefreshToken(newToken string) error
}

type TokensProvider struct {
	pp     repo.ParamsProvider
	tokens *tokens
}

func NewTokensProvider(pp repo.ParamsProvider) *TokensProvider {
	return &TokensProvider{
		pp: pp,
	}
}

func (p *TokensProvider) Init() error {
	if p.tokens != nil {
		return nil
	}

	tokens, err := p.loadTokens()
	if err != nil {
		return err
	}

	p.tokens = tokens
	return nil
}

func (p *TokensProvider) AccessToken() string {
	return p.tokens.AccessToken.Value
}

func (p *TokensProvider) UpdateAccessToken(newToken string) error {
	return p.updateToken(p.tokens.AccessToken, newToken, accessTokenPath)
}

func (p *TokensProvider) RefreshToken() string {
	return p.tokens.RefreshToken.Value
}

func (p *TokensProvider) UpdateRefreshToken(newToken string) error {
	return p.updateToken(p.tokens.RefreshToken, newToken, refreshTokenPath)
}

func (p *TokensProvider) loadTokens() (*tokens, error) {
	result := &tokens{
		AccessToken: &token{
			Value:     p.pp.GetParam(accessTokenPath),
			Updatable: true,
		},
		RefreshToken: &token{
			Value:     p.pp.GetParam(refreshTokenPath),
			Updatable: true,
		},
	}

	if err := result.validate(); err != nil {
		return nil, err
	}

	return result, nil
}

func (p *TokensProvider) updateToken(existingToken *token, newToken string, path string) error {
	if !existingToken.Updatable {
		return fmt.Errorf("token at path `%s` is not updatable", path)
	}

	if err := p.pp.UpdateParam(path, newToken); err != nil {
		return err
	}

	existingToken.Value = newToken
	return nil
}
