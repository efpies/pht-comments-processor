package lambda

import "pht/comments-processor/model"

type Event struct {
	ServiceRequest
	Platform model.Platform `json:"platform"`
}

type EventHandler interface {
	Handle(*Event) (any, error)
}
