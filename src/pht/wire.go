//go:build wireinject

package pht

import (
	"github.com/google/wire"
	"pht/comments-processor/google"
	"pht/comments-processor/pht/auth"
	"pht/comments-processor/pht/config"
	"pht/comments-processor/pht/handlers"
	"pht/comments-processor/pht/services"
	"pht/comments-processor/pht/sheets"
	"pht/comments-processor/pht/strategies"
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
	services.NewClient,

	wire.Bind(new(config.ConfigProvider), new(*config.Config)),

	providePostsProvider,

	wire.Bind(new(services.FixedPostsGetter), new(*services.PostsProvider)),
	wire.Bind(new(services.PostGetter), new(*services.PostsProvider)),
	wire.Bind(new(services.PostCommentsGetter), new(*services.Client)),
	wire.Bind(new(services.PagesGetter), new(*services.Client)),
	wire.Bind(new(services.WikiGetter), new(*services.Client)),

	google.NewConfig,
	wire.Bind(new(google.SheetsConfigProvider), new(*google.Config)),

	google.NewSheetsClient,
	sheets.NewDataProvider,
	sheets.NewGetPostsInfoStrategy,

	strategies.NewContentCheckPostStrategy,
	wire.Bind(new(strategies.CheckPostStrategy), new(*strategies.ContentCheckPostStrategy)),

	sheets.NewNotifierDataGetter,
)

func ProvideLocator(pp repo.ParamsProvider) (*locator, error) {
	wire.Build(PhtSet, newLocator)
	return nil, nil
}

func ProvideRouter(l *locator) (*handlers.Router, error) {
	wire.Build(
		handlers.NewRouter,
		wire.FieldsOf(new(*locator),
			"accessTokenProvider",
			"tokensRefresher",
			"fixedPostsGetter",
			"postGetter",
			"postCommentsGetter",
			"pagesGetter",
			"wikiGetter",
			"sheetsDataProvider",
			"checkPostStrategy",
			"getPostsInfoStrategy",
			"notifierDataGetter",
			"config",
		))
	return nil, nil
}

func providePostsProvider(c *services.Client) (*services.PostsProvider, error) {
	wire.Build(
		services.NewPostsProvider,
		wire.Bind(new(services.FixedPostsGetter), new(*services.Client)),
		wire.Bind(new(services.PostGetter), new(*services.Client)),
		wire.Bind(new(services.PostCommentsGetter), new(*services.Client)),
		wire.Bind(new(services.PagesGetter), new(*services.Client)),
		wire.Bind(new(services.WikiGetter), new(*services.Client)),
	)
	return nil, nil
}
