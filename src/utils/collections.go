package utils

import "github.com/samber/lo"

func TakeWhile[T any](collection []T, predicate func(item T) bool) []T {
	_, idx, ok := lo.FindIndexOf(collection, func(item T) bool {
		return !predicate(item)
	})
	if !ok {
		return collection
	}

	return collection[:idx]
}

func DropWhile[T any](collection []T, predicate func(item T) bool) []T {
	_, idx, ok := lo.FindIndexOf(collection, func(item T) bool {
		return !predicate(item)
	})
	if !ok {
		return collection
	}

	return collection[idx:]
}
