package model

import "errors"

type ErrorResponse struct {
	Detail string `json:"detail"`
	Code   string `json:"code"`
}

func (e *ErrorResponse) Error() error {
	return errors.New(e.Detail)
}
