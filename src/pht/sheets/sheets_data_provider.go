package sheets

import (
	"fmt"
	"github.com/samber/lo"
	"pht/comments-processor/google"
	"pht/comments-processor/pht/config"
)

type DataProvider struct {
	sheetsClient *google.SheetsClient
	config       config.ConfigProvider
}

func NewDataProvider(sheetsClient *google.SheetsClient, config config.ConfigProvider) *DataProvider {
	return &DataProvider{
		sheetsClient: sheetsClient,
		config:       config,
	}
}

func (provider *DataProvider) GetSheetData(flowID string, sheet string) ([][]string, error) {
	conf, ok := provider.config.FlowsSheets()[flowID]
	if !ok {
		return nil, fmt.Errorf("unknown flow ID: %s", flowID)
	}

	if !lo.Contains(conf.Sheets, sheet) {
		return nil, fmt.Errorf("unknown sheet: %s", sheet)
	}

	return provider.sheetsClient.GetSpreadsheetValues(conf.SpreadsheetID, sheet, "A1", "BZ250")
}
