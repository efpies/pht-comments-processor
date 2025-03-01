// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package pht

import (
	"github.com/google/wire"
	"pht/comments-processor/google"
	"pht/comments-processor/pht/auth"
	"pht/comments-processor/pht/config"
	"pht/comments-processor/pht/handlers"
	"pht/comments-processor/pht/services"
	"pht/comments-processor/repo"
)

// Injectors from wire.go:

func ProvideLocator(pp repo.ParamsProvider) (*Locator, error) {
	tokensProvider := auth.NewTokensProvider(pp)
	configConfig := config.NewConfig(pp)
	tokensRefresher, err := auth.NewTokensRefresher(configConfig, tokensProvider, tokensProvider, tokensProvider)
	if err != nil {
		return nil, err
	}
	client, err := services.NewClient(configConfig, tokensProvider, tokensRefresher)
	if err != nil {
		return nil, err
	}
	postsProvider, err := providePostsProvider(client)
	if err != nil {
		return nil, err
	}
	googleConfig := google.NewConfig(pp)
	sheetsClient, err := google.NewSheetsClient(googleConfig)
	if err != nil {
		return nil, err
	}
	sheetsDataProvider := services.NewSheetsDataProvider(sheetsClient, configConfig)
	locator := NewLocator(tokensProvider, tokensProvider, tokensRefresher, postsProvider, postsProvider, postsProvider, client, client, client, sheetsDataProvider)
	return locator, nil
}

func ProvideRouter(l *Locator) (*handlers.Router, error) {
	accessTokenProvider := l.accessTokenProvider
	tokensRefresher := l.tokensRefresher
	fixedPostsGetter := l.fixedPostsGetter
	postGetter := l.postGetter
	postCommentsGetter := l.postCommentsGetter
	pagesGetter := l.pagesGetter
	wikiGetter := l.wikiGetter
	sheetsDataProvider := l.sheetsDataProvider
	router := handlers.NewRouter(accessTokenProvider, tokensRefresher, fixedPostsGetter, postGetter, postCommentsGetter, pagesGetter, wikiGetter, sheetsDataProvider)
	return router, nil
}

func providePostsProvider(c *services.Client) (*services.PostsProvider, error) {
	postsProvider := services.NewPostsProvider(c, c, c, c, c)
	return postsProvider, nil
}

// wire.go:

var TokensProviderSet = wire.NewSet(auth.NewTokensProvider, wire.Bind(new(auth.AccessTokenProvider), new(*auth.TokensProvider)), wire.Bind(new(auth.AccessTokenUpdater), new(*auth.TokensProvider)), wire.Bind(new(auth.RefreshTokenProvider), new(*auth.TokensProvider)), wire.Bind(new(auth.RefreshTokenUpdater), new(*auth.TokensProvider)))

var PhtSet = wire.NewSet(
	TokensProviderSet, config.NewConfig, auth.NewTokensRefresher, services.NewClient, wire.Bind(new(config.ConfigProvider), new(*config.Config)), providePostsProvider, wire.Bind(new(services.FixedPostsGetter), new(*services.PostsProvider)), wire.Bind(new(services.PostGetter), new(*services.PostsProvider)), wire.Bind(new(services.PostCommentsGetter), new(*services.Client)), wire.Bind(new(services.PagesGetter), new(*services.Client)), wire.Bind(new(services.WikiGetter), new(*services.Client)), google.NewConfig, wire.Bind(new(google.SheetsConfigProvider), new(*google.Config)), google.NewSheetsClient, services.NewSheetsDataProvider,
)
