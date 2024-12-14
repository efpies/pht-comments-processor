package handlers

import (
	"pht/comments-processor/pht/auth"
)

type refreshAccessTokenDto struct {
	Access string `json:"access"`
}

func refreshAccessToken(tokensRefresher auth.TokensRefresher) lambdaHandlerOut[refreshAccessTokenDto] {
	return func() (refreshAccessTokenDto, error) {
		token, err := tokensRefresher.RefreshTokens()
		if err != nil {
			return refreshAccessTokenDto{}, err
		}

		return refreshAccessTokenDto{
			Access: token,
		}, nil
	}
}
