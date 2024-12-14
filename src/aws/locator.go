package aws

import (
	"context"
	"errors"
	sdkConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"pht/comments-processor/repo"
	"pht/comments-processor/services"
)

type locator struct {
	paramsProvider repo.ParamsProvider
}

func NewLocator(paramsPath string) (services.InfraLocator, error) {
	cfg, err := sdkConfig.LoadDefaultConfig(
		context.TODO(),
		sdkConfig.WithSharedConfigProfile(NewConfig().Profile()))
	if err != nil {
		return nil, errors.Join(errors.New("couldn't create AWS SDK services"), err)
	}

	return &locator{
		paramsProvider: NewSSMParamsProvider(ssm.NewFromConfig(cfg), paramsPath),
	}, nil
}

func (l *locator) Init() error {
	if err := l.paramsProvider.PrefetchParams(); err != nil {
		return errors.Join(errors.New("couldn't init params provider"), err)
	}

	return nil
}

func (l *locator) ParamsProvider() repo.ParamsProvider {
	return l.paramsProvider
}
