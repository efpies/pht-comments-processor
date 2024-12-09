package lambda

type ServiceRequest struct {
	Method string         `json:"method"`
	Params map[string]any `json:"params"`
}

type ServiceRequestHandler interface {
	Handle(*ServiceRequest) (any, error)
}
