package google

import "pht/comments-processor/repo"

type SheetsConfigProvider interface {
	SheetsServiceKey() string
}

type Config struct {
	pp repo.ParamsProvider
}

func NewConfig(pp repo.ParamsProvider) *Config {
	return &Config{
		pp: pp,
	}
}

func (c *Config) SheetsServiceKey() string {
	return c.pp.GetParam("google/sheetsServiceKey")
}
