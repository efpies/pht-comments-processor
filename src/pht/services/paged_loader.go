package services

import (
	"pht/comments-processor/pht/model"
	"slices"
)

func PagedLoad[T any](pageFrom int, pageToInclusive int, loader func(page int) (model.Page[T], error)) ([]T, error) {
	var response []T

	page := pageFrom

	for {
		if page > pageToInclusive {
			break
		}

		subResponse, err := loader(page)
		if err != nil {
			return []T{}, err
		}

		response = slices.Concat(response, subResponse.Items)

		if !subResponse.Pagination.HasNextPage {
			break
		}

		page++
	}

	return response, nil
}
