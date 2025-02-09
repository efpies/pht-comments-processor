package model

import "errors"

type ErrorResponse struct {
	Detail string `json:"detail"`
	Code   string `json:"code"`
}

type ErrorDto struct {
	Error string `json:"error"`
}

func (e *ErrorResponse) Error() error {
	return errors.New(e.Detail)
}
