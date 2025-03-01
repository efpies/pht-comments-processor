package config

import (
	"encoding/json"
	"pht/comments-processor/repo"
)

type ConfigProvider interface {
	RefreshTokenURL() string
	ContentURL() string
	FlowsSheets() map[string]FlowSheetConfig
}

type Config struct {
	pp repo.ParamsProvider
}

type FlowSheetConfig struct {
	SpreadsheetID string   `json:"spreadsheetId"`
	Sheets        []string `json:"sheets"`
}

func NewConfig(pp repo.ParamsProvider) *Config {
	return &Config{
		pp: pp,
	}
}

func (c Config) RefreshTokenURL() string {
	return c.pp.GetParam("pht/refreshTokenUrl")
}

func (c Config) ContentURL() string {
	return c.pp.GetParam("pht/contentUrl")
}

func (c Config) FlowsSheets() map[string]FlowSheetConfig {
	p := c.pp.GetParam("pht/flows/sheets")

	var flows map[string]FlowSheetConfig
	if err := json.Unmarshal([]byte(p), &flows); err != nil {
		panic(err)
	}

	return flows
}
