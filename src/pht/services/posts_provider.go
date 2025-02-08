package services

import (
	"context"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"
	"log"
	"pht/comments-processor/pht/model"
	"pht/comments-processor/utils"
	"strconv"
	"sync"
)

const (
	feedListKey  string = "feed"
	wikiListKey  string = "wiki"
	topicListKey string = "topic"
)

type cursor struct {
	curPage int
	hasMore bool
}

type PostsProvider struct {
	cache         sync.Map
	fixedPostsIDs []int
	cursors       map[string]*cursor

	postGetter       PostGetter
	fixedPostsGetter FixedPostsGetter
	pagesGetter      PagesGetter
	wikiGetter       WikiGetter
	postFiller       *PostFiller

	cursorMutex   sync.Mutex
	pageLoadTasks sync.Map

	init func() error
}

func NewPostsProvider(
	postGetter PostGetter,
	fixedPostsGetter FixedPostsGetter,
	pagesGetter PagesGetter,
	wikiGetter WikiGetter,
	postCommentsGetter PostCommentsGetter) *PostsProvider {
	p := &PostsProvider{
		cache:            sync.Map{},
		fixedPostsIDs:    []int{},
		cursors:          map[string]*cursor{},
		pageLoadTasks:    sync.Map{},
		postGetter:       postGetter,
		fixedPostsGetter: fixedPostsGetter,
		pagesGetter:      pagesGetter,
		wikiGetter:       wikiGetter,
		postFiller:       NewPostFiller(postCommentsGetter),
	}
	p.init = sync.OnceValue(p.doInit)
	return p
}

func (p *PostsProvider) Init() error {
	return p.init()
}

func (p *PostsProvider) doInit() error {
	var err error
	var eg errgroup.Group
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	eg.Go(func() error { return p.loadWikis(&eg, ctx, cancel) })
	eg.Go(func() error { return p.loadFixedPosts(ctx, cancel) })

	eg.Go(func() error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		_, err = p.preloadNextPage(feedListKey, nil, false)
		return err
	})
	eg.Go(func() error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		_, err = p.preloadNextPage(topicListKey, nil, false)
		return err
	})

	return eg.Wait()
}

func (p *PostsProvider) getCachedPost(id int) (model.PostDto, bool) {
	value, ok := p.cache.Load(id)
	if !ok {
		return model.PostDto{}, false
	}

	return value.(model.PostDto), ok
}

func (p *PostsProvider) cachePost(id int, post model.PostDto) {
	p.cache.Store(id, post)
}

func (p *PostsProvider) loadFixedPosts(ctx context.Context, cancel context.CancelFunc) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	fixedPosts, err := p.fixedPostsGetter.GetFixedPosts()
	if err != nil {
		cancel()
		return err
	}

	for _, post := range fixedPosts {
		p.cachePost(post.ID, post)
	}

	p.fixedPostsIDs = lo.Map(fixedPosts, func(post model.PostDto, _ int) int {
		return post.ID
	})
	return nil
}

func (p *PostsProvider) loadWikis(eg *errgroup.Group, ctx context.Context, cancel context.CancelFunc) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	wikis, err := p.wikiGetter.GetWikis()
	if err != nil {
		cancel()
		return err
	}

	for _, wiki := range wikis {
		if wiki.ID == 6 {
			// Recipes
			continue
		}

		eg.Go(func() error { return p.loadWiki(wiki.ID, ctx, cancel) })
	}

	return nil
}

func (p *PostsProvider) loadWiki(id int, ctx context.Context, cancel context.CancelFunc) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		hasNext, err := p.preloadNextPage(wikiListKey, utils.Ptr(strconv.Itoa(id)), true)
		if err != nil {
			cancel()
			return err
		}
		if !hasNext {
			break
		}
	}
	return nil
}

func (p *PostsProvider) GetFixedPosts() ([]model.PostDto, error) {
	posts := make([]model.PostDto, 0, len(p.fixedPostsIDs))
	for _, postID := range p.fixedPostsIDs {
		post, err := p.GetPost(postID)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (p *PostsProvider) GetPost(postID int) (model.PostDto, error) {
	var post model.PostDto
	var err error
	var ok bool

	post, ok, err = p.tryGetCachedPost(postID, feedListKey, nil)
	if err != nil {
		return model.PostDto{}, err
	}

	if !ok {
		post, err = p.postGetter.GetPost(postID)
		if err != nil {
			return model.PostDto{}, err
		}

		p.cachePost(postID, post)
	}

	if err := p.postFiller.FillPost(&post); err != nil {
		return model.PostDto{}, err
	}

	p.cachePost(postID, post)
	return post, nil
}

func (p *PostsProvider) tryGetCachedPost(postID int, list string, sublist *string) (post model.PostDto, ok bool, err error) {
	for i := 0; i < 3; i++ {
		post, ok = p.getCachedPost(postID)
		if ok {
			return post, true, nil
		}

		hasNext, err := p.preloadNextPage(list, sublist, false)
		if err != nil {
			return model.PostDto{}, false, err
		}

		if !hasNext {
			break
		}
	}

	return model.PostDto{}, false, nil
}

func (p *PostsProvider) preloadNextPage(list string, sublist *string, loadAll bool) (hasNext bool, err error) {
	cacheKey := list
	if sublist != nil {
		cacheKey += "/" + *sublist
	}

	p.cursorMutex.Lock()
	listCursor, ok := p.cursors[cacheKey]
	if !ok {
		listCursor = &cursor{curPage: 1, hasMore: true}
		p.cursors[cacheKey] = listCursor
	}
	p.cursorMutex.Unlock()

	ch, ok := p.pageLoadTasks.LoadOrStore(cacheKey, make(chan error, 1))
	if ok {
		err = <-ch.(chan error)
		return listCursor.hasMore, err
	}
	defer p.pageLoadTasks.Delete(cacheKey)
	defer close(ch.(chan error))

	if !listCursor.hasMore {
		return false, nil
	}

	log.Printf("[%s p%d] Loading (all: %t)", cacheKey, listCursor.curPage, loadAll)
	to := &listCursor.curPage
	if loadAll {
		to = nil
	}

	posts, hasMore, err := p.pagesGetter.GetPages(listCursor.curPage, to, list, sublist)
	if err != nil {
		log.Printf("[%s p%d] Couldn't load (all: %t): %v", cacheKey, listCursor.curPage, loadAll, err)
		ch.(chan error) <- err
		return false, err
	}

	for _, post := range posts {
		p.cachePost(post.ID, post)
	}

	if !hasMore {
		log.Printf("[%s p%d] Last page (all: %t)", cacheKey, listCursor.curPage, loadAll)
		listCursor.hasMore = false
	} else {
		listCursor.curPage++
	}

	return hasMore, nil
}
