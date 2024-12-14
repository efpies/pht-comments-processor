package handlers

import (
	"pht/comments-processor/pht/auth"
)

type getAccessTokenDto struct {
	Access string `json:"access"`
}

func getAccessToken(tokenProvider auth.AccessTokenProvider) lambdaHandlerOut[getAccessTokenDto] {
	return func() (getAccessTokenDto, error) {
		return getAccessTokenDto{
			Access: tokenProvider.AccessToken(),
		}, nil
	}
}
