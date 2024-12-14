package pht

import (
	"errors"
	"pht/comments-processor/pht/auth"
	"pht/comments-processor/pht/config"
	"pht/comments-processor/repo"
)

type Locator interface {
	Config() config.ConfigProvider
	AccessTokenProvider() auth.AccessTokenProvider
	AccessTokenUpdater() auth.AccessTokenUpdater
	RefreshTokenProvider() auth.RefreshTokenProvider
	RefreshTokenUpdater() auth.RefreshTokenUpdater
	TokensRefresher() auth.TokensRefresher
}

type locator struct {
	config          *config.Config
	tokensProvider  *auth.TokensProvider
	tokensRefresher auth.TokensRefresher
}

func NewLocator(pp repo.ParamsProvider) (Locator, error) {
	l := &locator{
		config:         config.NewConfig(pp),
		tokensProvider: auth.NewTokensProvider(pp),
	}

	if err := l.init(); err != nil {
		return nil, err
	}

	tr, err := auth.NewTokensRefresher(l.Config(), l.RefreshTokenProvider(), l.AccessTokenUpdater(), l.RefreshTokenUpdater())
	if err != nil {
		return nil, err
	}
	l.tokensRefresher = tr

	return l, nil
}

func (l *locator) init() error {
	if err := l.tokensProvider.Init(); err != nil {
		return errors.Join(errors.New("couldn't init tokens provider"), err)
	}

	return nil
}

func (l *locator) Config() config.ConfigProvider {
	return l.config
}

func (l *locator) AccessTokenProvider() auth.AccessTokenProvider {
	return l.tokensProvider
}

func (l *locator) AccessTokenUpdater() auth.AccessTokenUpdater {
	return l.tokensProvider
}

func (l *locator) RefreshTokenProvider() auth.RefreshTokenProvider {
	return l.tokensProvider
}

func (l *locator) RefreshTokenUpdater() auth.RefreshTokenUpdater {
	return l.tokensProvider
}

func (l *locator) TokensRefresher() auth.TokensRefresher {
	return l.tokensRefresher
}
