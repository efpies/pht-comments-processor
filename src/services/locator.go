package services

import (
	"pht/comments-processor/repo"
)

type InfraLocator interface {
	Init() error
	ParamsProvider() repo.ParamsProvider
}
