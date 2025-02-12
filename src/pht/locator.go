package pht

import (
	"errors"
	"pht/comments-processor/pht/auth"
	"pht/comments-processor/pht/services"
)

type Locator struct {
	tokensProvider      *auth.TokensProvider
	accessTokenProvider auth.AccessTokenProvider
	tokensRefresher     auth.TokensRefresher
	postsProvider       *services.PostsProvider
	fixedPostsGetter    services.FixedPostsGetter
	postGetter          services.PostGetter
	postCommentsGetter  services.PostCommentsGetter
	pagesGetter         services.PagesGetter
	wikiGetter          services.WikiGetter
}

func NewLocator(
	tokensProvider *auth.TokensProvider,
	accessTokenProvider auth.AccessTokenProvider,
	tokensRefresher auth.TokensRefresher,
	postsProvider *services.PostsProvider,
	fixedPostsGetter services.FixedPostsGetter,
	postGetter services.PostGetter,
	postCommentsGetter services.PostCommentsGetter,
	pagesGetter services.PagesGetter,
	wikiGetter services.WikiGetter,
) *Locator {
	return &Locator{
		tokensProvider:      tokensProvider,
		accessTokenProvider: accessTokenProvider,
		tokensRefresher:     tokensRefresher,
		postsProvider:       postsProvider,
		fixedPostsGetter:    fixedPostsGetter,
		postGetter:          postGetter,
		postCommentsGetter:  postCommentsGetter,
		pagesGetter:         pagesGetter,
		wikiGetter:          wikiGetter,
	}
}

func (s *Locator) Init() error {
	if err := s.tokensProvider.Init(); err != nil {
		return errors.Join(errors.New("couldn't init tokens provider"), err)
	}

	if err := s.postsProvider.Init(); err != nil {
		return errors.Join(errors.New("couldn't init posts provider"), err)
	}

	return nil
}
