package pht

import (
	"errors"
	"pht/comments-processor/pht/auth"
)

type Locator struct {
	tokensProvider      *auth.TokensProvider
	accessTokenProvider auth.AccessTokenProvider
	tokensRefresher     auth.TokensRefresher
}

func NewLocator(
	tokensProvider *auth.TokensProvider,
	accessTokenProvider auth.AccessTokenProvider,
	tokensRefresher auth.TokensRefresher,
) *Locator {
	return &Locator{
		tokensProvider:      tokensProvider,
		accessTokenProvider: accessTokenProvider,
		tokensRefresher:     tokensRefresher,
	}
}

func (s *Locator) Init() error {
	if err := s.tokensProvider.Init(); err != nil {
		return errors.Join(errors.New("couldn't init tokens provider"), err)
	}

	return nil
}
