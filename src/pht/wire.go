//go:build wireinject

package pht

import (
	"github.com/google/wire"
	"pht/comments-processor/pht/auth"
	"pht/comments-processor/pht/config"
	"pht/comments-processor/pht/handlers"
	"pht/comments-processor/repo"
)

var TokensProviderSet = wire.NewSet(
	auth.NewTokensProvider,
	wire.Bind(new(auth.AccessTokenProvider), new(*auth.TokensProvider)),
	wire.Bind(new(auth.AccessTokenUpdater), new(*auth.TokensProvider)),
	wire.Bind(new(auth.RefreshTokenProvider), new(*auth.TokensProvider)),
	wire.Bind(new(auth.RefreshTokenUpdater), new(*auth.TokensProvider)),
)

var PhtSet = wire.NewSet(
	TokensProviderSet,

	config.NewConfig,
	auth.NewTokensRefresher,

	wire.Bind(new(config.ConfigProvider), new(*config.Config)),

	NewLocator,
)

func ProvideLocator(pp repo.ParamsProvider) (*Locator, error) {
	wire.Build(PhtSet)
	return nil, nil
}

func ProvideRouter(l *Locator) (*handlers.Router, error) {
	wire.Build(
		handlers.NewRouter,
		wire.FieldsOf(new(*Locator),
			"accessTokenProvider",
			"tokensRefresher",
		))
	return nil, nil
}
