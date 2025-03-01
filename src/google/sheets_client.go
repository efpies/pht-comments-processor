package google

import (
	"context"
	"fmt"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type SheetsClient struct {
	*sheets.Service
}

func NewSheetsClient(config SheetsConfigProvider) (*SheetsClient, error) {
	srv, err := sheets.NewService(context.TODO(), option.WithAPIKey(config.SheetsServiceKey()))
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Sheets client: %v", err)
	}

	return &SheetsClient{
		Service: srv,
	}, nil
}

func (c *SheetsClient) GetSpreadsheetValues(spreadsheetID string, sheet string, rangeFrom string, rangeTo string) ([][]string, error) {
	resp, err := c.Spreadsheets.Values.Get(spreadsheetID, fmt.Sprintf("%s!%s:%s", sheet, rangeFrom, rangeTo)).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve data from sheet: %v", err)
	}

	if len(resp.Values) == 0 {
		return nil, fmt.Errorf("no data found")
	}

	var result [][]string
	for _, row := range resp.Values {
		stringRow := make([]string, 0)
		for _, cell := range row {
			stringRow = append(stringRow, fmt.Sprintf("%v", cell))
		}
		result = append(result, stringRow)
	}

	return result, nil
}
