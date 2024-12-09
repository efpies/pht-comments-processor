package services

import (
	"fmt"
	"pht/comments-processor/handlers/lambda"
	"pht/comments-processor/model"
)

type LambdaEventHandler struct {
	platformHandlers map[model.Platform]lambda.ServiceRequestHandler
}

func NewLambdaEventHandler() *LambdaEventHandler {
	return &LambdaEventHandler{
		platformHandlers: map[model.Platform]lambda.ServiceRequestHandler{},
	}
}

func (h *LambdaEventHandler) RegisterPlatformHandler(platform model.Platform, handler lambda.ServiceRequestHandler) error {
	if handler == nil {
		return fmt.Errorf("handler must not be nil")
	}

	h.platformHandlers[platform] = handler
	return nil
}

func (h *LambdaEventHandler) Handle(event *lambda.Event) (any, error) {
	if event == nil {
		return nil, fmt.Errorf("received nil event")
	}

	handler := h.platformHandlers[event.Platform]
	if handler == nil {
		return nil, fmt.Errorf("unhandled platform type: %s", event.Platform)
	}

	return handler.Handle(&event.ServiceRequest)
}
