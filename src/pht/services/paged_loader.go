package services

import (
	"pht/comments-processor/pht/model"
	"slices"
)

func PagedLoad[T any](pageFrom int, pageToInclusive int, loader func(page int) (model.Page[T], error)) (items []T, hasMore bool, err error) {
	page := pageFrom
	hasMore = true

	for {
		if page > pageToInclusive {
			break
		}

		subResponse, err := loader(page)
		if err != nil {
			return []T{}, false, err
		}

		items = slices.Concat(items, subResponse.Items)

		if !subResponse.Pagination.HasNextPage {
			hasMore = false
			break
		}

		page++
	}

	return items, hasMore, nil
}
