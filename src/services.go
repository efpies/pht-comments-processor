package main

import (
	"pht/comments-processor/aws"
	"pht/comments-processor/services"
)

type appServices struct {
	infraServices services.InfraLocator
}

func newAppServices() (*appServices, error) {
	s, err := aws.NewLocator()
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

	return nil
}
