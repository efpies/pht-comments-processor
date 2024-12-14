package main

import (
	"pht/comments-processor/aws"
	"pht/comments-processor/handlers/lambda"
	"pht/comments-processor/model"
	"pht/comments-processor/pht"
	"pht/comments-processor/services"
)

type appServices struct {
	infraServices      services.InfraLocator
	phtServices        *pht.Locator
	lambdaEventHandler lambda.EventHandler
}

func newAppServices() (*appServices, error) {
	s, err := aws.NewLocator("/pht-comments-processor")
	if err != nil {
		return nil, err
	}

	return &appServices{
		infraServices: s,
	}, nil
}

func (s *appServices) init() error {
	if err := s.infraServices.Init(); err != nil {
		return err
	}

	phtServices, err := pht.ProvideLocator(s.infraServices.ParamsProvider())
	if err != nil {
		return err
	}

	if err = phtServices.Init(); err != nil {
		return err
	}

	phtRouter, err := pht.ProvideRouter(phtServices)
	if err != nil {
		return err
	}

	leh := services.NewLambdaEventHandler()
	if err = leh.RegisterPlatformHandler(model.PlatformEnum.Pht, phtRouter); err != nil {
		return err
	}

	s.lambdaEventHandler = leh

	return nil
}
