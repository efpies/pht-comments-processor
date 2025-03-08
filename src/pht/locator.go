package pht

import (
	"errors"
	"pht/comments-processor/pht/auth"
	"pht/comments-processor/pht/config"
	"pht/comments-processor/pht/services"
	"pht/comments-processor/pht/sheets"
	"pht/comments-processor/pht/strategies"
	cservices "pht/comments-processor/services"
)

type Locator interface {
	cservices.Initable
	GetNotifierDataGetter() *sheets.NotifierDataGetter
}

type locator struct {
	tokensProvider       *auth.TokensProvider
	accessTokenProvider  auth.AccessTokenProvider
	tokensRefresher      auth.TokensRefresher
	postsProvider        *services.PostsProvider
	fixedPostsGetter     services.FixedPostsGetter
	postGetter           services.PostGetter
	postCommentsGetter   services.PostCommentsGetter
	pagesGetter          services.PagesGetter
	wikiGetter           services.WikiGetter
	sheetsDataProvider   *sheets.DataProvider
	checkPostStrategy    strategies.CheckPostStrategy
	getPostsInfoStrategy *sheets.GetPostsInfoStrategy
	notifierDataGetter   *sheets.NotifierDataGetter
	config               config.ConfigProvider
}

func newLocator(
	tokensProvider *auth.TokensProvider,
	accessTokenProvider auth.AccessTokenProvider,
	tokensRefresher auth.TokensRefresher,
	postsProvider *services.PostsProvider,
	fixedPostsGetter services.FixedPostsGetter,
	postGetter services.PostGetter,
	postCommentsGetter services.PostCommentsGetter,
	pagesGetter services.PagesGetter,
	wikiGetter services.WikiGetter,
	sheetsDataProvider *sheets.DataProvider,
	checkPostStrategy strategies.CheckPostStrategy,
	getPostsInfoStrategy *sheets.GetPostsInfoStrategy,
	notifierDataGetter *sheets.NotifierDataGetter,
	config config.ConfigProvider,
) *locator {
	return &locator{
		tokensProvider:       tokensProvider,
		accessTokenProvider:  accessTokenProvider,
		tokensRefresher:      tokensRefresher,
		postsProvider:        postsProvider,
		fixedPostsGetter:     fixedPostsGetter,
		postGetter:           postGetter,
		postCommentsGetter:   postCommentsGetter,
		pagesGetter:          pagesGetter,
		wikiGetter:           wikiGetter,
		sheetsDataProvider:   sheetsDataProvider,
		checkPostStrategy:    checkPostStrategy,
		getPostsInfoStrategy: getPostsInfoStrategy,
		notifierDataGetter:   notifierDataGetter,
		config:               config,
	}
}

func (s *locator) Init() error {
	if err := s.tokensProvider.Init(); err != nil {
		return errors.Join(errors.New("couldn't init tokens provider"), err)
	}

	if err := s.postsProvider.Init(); err != nil {
		return errors.Join(errors.New("couldn't init posts provider"), err)
	}

	return nil
}

func (s *locator) GetNotifierDataGetter() *sheets.NotifierDataGetter {
	return s.notifierDataGetter
}
