package pht

import (
	"errors"
	"pht/comments-processor/pht/auth"
	"pht/comments-processor/pht/config"
	"pht/comments-processor/pht/services"
)

type Locator struct {
	config              config.ConfigProvider
	tokensProvider      *auth.TokensProvider
	accessTokenProvider auth.AccessTokenProvider
	tokensRefresher     auth.TokensRefresher
	fixedPostsGetter    services.FixedPostsGetter
}

func NewLocator(
	config config.ConfigProvider,
	tokensProvider *auth.TokensProvider,
	accessTokenProvider auth.AccessTokenProvider,
	tokensRefresher auth.TokensRefresher,
	fixedPostsGetter services.FixedPostsGetter,
) *Locator {
	return &Locator{
		config:              config,
		tokensProvider:      tokensProvider,
		accessTokenProvider: accessTokenProvider,
		tokensRefresher:     tokensRefresher,
		fixedPostsGetter:    fixedPostsGetter,
	}
}

func (s *Locator) Init() error {
	if err := s.tokensProvider.Init(); err != nil {
		return errors.Join(errors.New("couldn't init tokens provider"), err)
	}

	return nil
}
