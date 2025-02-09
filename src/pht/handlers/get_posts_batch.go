package handlers

import (
	"pht/comments-processor/pht/model"
	"pht/comments-processor/pht/services"
	"sync"
)

type getPostsBatchRequest struct {
	PostIDs []int `mapstructure:"post_ids"`
}

func getPostsBatch(postGetter services.PostGetter) lambdaHandlerInOut[getPostsBatchRequest, map[int]any] {
	return func(req getPostsBatchRequest) (map[int]any, error) {
		var wg sync.WaitGroup
		var mu sync.Mutex

		wg.Add(len(req.PostIDs))
		posts := make(map[int]any, len(req.PostIDs))
		for _, postID := range req.PostIDs {
			go func(postID int) {
				defer wg.Done()
				post, err := postGetter.GetPost(postID)

				mu.Lock()
				defer mu.Unlock()
				if err == nil {
					posts[postID] = post
				} else {
					posts[postID] = model.ErrorDto{
						Error: err.Error(),
					}
				}
			}(postID)
		}

		wg.Wait()

		return posts, nil
	}
}
