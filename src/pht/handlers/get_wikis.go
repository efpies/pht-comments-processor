package handlers

import (
	"pht/comments-processor/pht/model"
	"pht/comments-processor/pht/services"
)

func getWikis(wikiGetter services.WikiGetter) lambdaHandlerOut[[]model.WikiDto] {
	return func() ([]model.WikiDto, error) {
		return wikiGetter.GetWikis()
	}
}
