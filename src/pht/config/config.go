package config

import "pht/comments-processor/repo"

type ConfigProvider interface {
	RefreshTokenURL() string
	ContentURL() string
}

type Config struct {
	pp repo.ParamsProvider
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
