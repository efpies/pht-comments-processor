package aws

import (
	"context"
	"errors"
	sdkConfig "github.com/aws/aws-sdk-go-v2/config"
	"pht/comments-processor/services"
)

type locator struct {
}

func NewLocator() (services.InfraLocator, error) {
	_, err := sdkConfig.LoadDefaultConfig(
		context.TODO(),
		sdkConfig.WithSharedConfigProfile(NewConfig().Profile()))
	if err != nil {
		return nil, errors.Join(errors.New("couldn't create AWS SDK services"), err)
	}

	s := &locator{}
	return s, nil
}

func (l *locator) Init() error {
	return nil
}
