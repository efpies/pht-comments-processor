package handlers

import (
	"pht/comments-processor/pht/model/dto"
	"pht/comments-processor/pht/sheets"
)

type getNotifierDataRequest struct {
	PostType sheets.PostType `mapstructure:"post_type"`
	FlowID   string          `mapstructure:"flow_id"`
	Sheet    string          `mapstructure:"sheet"`
}

func getNotifierData(notifierDataGetter *sheets.NotifierDataGetter) lambdaHandlerInOut[getNotifierDataRequest, map[string]dto.NotifierData] {
	return func(req getNotifierDataRequest) (map[string]dto.NotifierData, error) {
		return notifierDataGetter.GetNotifierData(req.PostType, req.FlowID, req.Sheet)
	}
}
