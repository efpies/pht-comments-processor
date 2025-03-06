package handlers

import (
	"pht/comments-processor/pht/model/dto"
	"pht/comments-processor/pht/services"
)

func getWikis(wikiGetter services.WikiGetter) lambdaHandlerOut[[]dto.Wiki] {
	return func() ([]dto.Wiki, error) {
		return wikiGetter.GetWikis()
	}
}
