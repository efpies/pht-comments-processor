package handlers

import (
	"pht/comments-processor/pht/services"
)

type getSheetDataRequest struct {
	FlowID string `mapstructure:"flow_id"`
	Sheet  string `mapstructure:"sheet"`
}

func getSheetData(sheetsDataProvider *services.SheetsDataProvider) lambdaHandlerInOut[getSheetDataRequest, [][]string] {
	return func(req getSheetDataRequest) ([][]string, error) {
		return sheetsDataProvider.GetSheetData(req.FlowID, req.Sheet)
	}
}
