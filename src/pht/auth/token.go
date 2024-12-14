package auth

import "errors"

type token struct {
	Value     string
	Updatable bool
}

type tokens struct {
	AccessToken  *token
	RefreshToken *token
}

func (t *tokens) validate() error {
	if t.RefreshToken == nil {
		return errors.New("refresh token is missing")
	}

	return nil
}
