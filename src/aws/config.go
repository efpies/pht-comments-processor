package aws

import "os"

type ConfigProvider interface {
	Profile() string
}

type Config struct {
}

func NewConfig() *Config {
	return &Config{}
}

func (_ Config) Profile() string {
	return os.Getenv("AWS_PROFILE")
}
