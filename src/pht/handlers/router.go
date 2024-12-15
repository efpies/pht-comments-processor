package handlers

import (
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"pht/comments-processor/handlers/lambda"
	"pht/comments-processor/pht/auth"
	"pht/comments-processor/pht/services"
	"reflect"
)

type lambdaHandlerOut[TResp any] = func() (TResp, error)
type lambdaHandlerInOut[TReq any, TResp any] = func(TReq) (TResp, error)

type Router struct {
	accessTokenProvider auth.AccessTokenProvider
	tokensRefresher     auth.TokensRefresher
	fixedPostsGetter    services.FixedPostsGetter
	postGetter          services.PostGetter
	postCommentsGetter  services.PostCommentsGetter
	wikiGetter          services.WikiGetter
}

func NewRouter(
	accessTokenProvider auth.AccessTokenProvider,
	tokensRefresher auth.TokensRefresher,
	fixedPostsGetter services.FixedPostsGetter,
	postGetter services.PostGetter,
	postCommentsGetter services.PostCommentsGetter,
	wikiGetter services.WikiGetter,
) *Router {
	return &Router{
		accessTokenProvider: accessTokenProvider,
		tokensRefresher:     tokensRefresher,
		fixedPostsGetter:    fixedPostsGetter,
		postGetter:          postGetter,
		postCommentsGetter:  postCommentsGetter,
		wikiGetter:          wikiGetter,
	}
}

func (r *Router) Handle(request *lambda.ServiceRequest) (any, error) {
	if request == nil {
		return nil, fmt.Errorf("received nil request")
	}

	handler, err := r.makeHandler(request.Method)
	if err != nil {
		return nil, err
	}

	handlerValue := reflect.ValueOf(handler)
	handlerType := handlerValue.Type()
	if handlerType.Kind() != reflect.Func {
		return nil, fmt.Errorf("handler must be a function")
	}
	if handlerType.NumOut() != 2 {
		return nil, fmt.Errorf("handler must return 2 values: result and error")
	}

	var args []reflect.Value
	switch handlerType.NumIn() {
	case 0:
		break
	case 1:
		req := reflect.New(handlerType.In(0))
		if err = mapstructure.Decode(request.Params, req.Interface()); err != nil {
			return nil, errors.Join(errors.New("failed to decode request"), err)
		}

		args = append(args, req.Elem())
		break
	default:
		return nil, fmt.Errorf("handler must accept at most 1 argument")
	}

	results := handlerValue.Call(args)
	err, _ = results[1].Interface().(error)

	return results[0].Interface(), err
}

func (r *Router) makeHandler(method string) (any, error) {
	switch method {
	case "token/access":
		return getAccessToken(r.accessTokenProvider), nil
	case "token/refresh":
		return refreshAccessToken(r.tokensRefresher), nil
	case "content/post/fixed":
		return getFixedPosts(r.fixedPostsGetter, r.postCommentsGetter), nil
	case "content/post/by-id":
		return getPost(r.postGetter, r.postCommentsGetter), nil
	case "content/post/comments/list":
		return getPostComments(r.postCommentsGetter), nil
	case "content/wiki/list":
		return getWikis(r.wikiGetter), nil
	default:
		return nil, fmt.Errorf("unhandled method: %s", method)
	}
}
