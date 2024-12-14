package services

import (
	"pht/comments-processor/repo"
)

type InfraLocator interface {
	Initable
	ParamsProvider() repo.ParamsProvider
}
