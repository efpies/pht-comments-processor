package lambda

import "pht/comments-processor/model"

type Event struct {
	ServiceRequest
	Platform model.Platform `json:"platform"`
}

type RouteEvent struct {
	RawPath string `json:"rawPath"`
}

type EventHandler interface {
	Handle(*Event) (any, error)
}
