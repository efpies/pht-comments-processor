package http

type emptyURLError struct{}

func (e *emptyURLError) Error() string {
	return "base url is empty"
}

type noRequestSenderError struct{}

func (e *noRequestSenderError) Error() string {
	return "request sender is not provided"
}

type noSerializerError struct{}

func (e *noSerializerError) Error() string {
	return "serializer is not provided"
}
