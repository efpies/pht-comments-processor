package auth

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/awsdocs/aws-doc-sdk-examples/gov2/testtools"
	"pht/comments-processor/aws"
	"pht/comments-processor/utils"
	"testing"
)

const baseParamPath = "/path"

func TestTokensCacheUsed(t *testing.T) {
	t.Run("Tokens cache is used", func(t *testing.T) {
		provider := TokensProvider{
			nil, // don't load tokens from the store
			&tokens{
				AccessToken: &token{
					Value:     "access_tok",
					Updatable: true,
				},
				RefreshToken: &token{
					Value:     "refresh_tok",
					Updatable: true,
				},
			}}

		tokens := provider.tokens
		if tokens == nil {
			t.Errorf("couldn't get cached tokens")
		}

		if err := tokens.validate(); err != nil {
			t.Errorf("invalid cached tokens: %v", err)
		}
	})
}

func TestTokensReadCorrectly(t *testing.T) {
	correctAccessToken := "access_tok"
	correctRefreshToken := "refresh_tok"

	cases := []struct {
		accessToken  *string
		refreshToken *string
	}{
		{
			accessToken:  &correctAccessToken,
			refreshToken: &correctRefreshToken,
		},
		{
			refreshToken: &correctRefreshToken,
		},
		{
			accessToken: &correctAccessToken,
		},
	}

	for _, tt := range cases {
		t.Run("Tokens read correctly", func(t *testing.T) {
			stubber := testtools.NewStubber()

			var params []types.Parameter
			if tt.accessToken != nil {
				params = append(params, makeParam(tt.accessToken, baseParamPath+"/"+accessTokenPath))
			}
			if tt.refreshToken != nil {
				params = append(params, makeParam(tt.refreshToken, baseParamPath+"/"+refreshTokenPath))
			}

			stubber.Add(makeGetParametersByPathStub(params))

			provider := setupProvider(t, stubber)
			if err := provider.Init(); err != nil {
				t.Fatalf("couldn't load tokens: %s", err)
			}

			tokens := provider.tokens
			if tokens == nil {
				t.Errorf("tokens are missing")
			} else {
				if tt.accessToken != nil && (tokens.AccessToken == nil || tokens.AccessToken.Value != *tt.accessToken) {
					t.Errorf("expected access token to be `%v` but got `%v`", *tt.accessToken, tokens.AccessToken.Value)
				}
				if tt.refreshToken != nil && (tokens.RefreshToken == nil || tokens.RefreshToken.Value != *tt.refreshToken) {
					t.Errorf("expected refresh token to be `%v` but got `%v`", *tt.refreshToken, tokens.RefreshToken.Value)
				}
			}
		})
	}
}

func TestTokensUpdatedCorrectly(t *testing.T) {
	accessToken := "access_tok"
	refreshToken := "refresh_tok"

	cases := []struct {
		accessTokenUpdated bool
		wantErr            error
	}{
		{
			accessTokenUpdated: true,
		},
		{
			accessTokenUpdated: false,
			wantErr:            fmt.Errorf("token at path `%v` is not updatable", accessTokenPath),
		},
	}

	for _, tt := range cases {
		t.Run("Tokens updated correctly", func(t *testing.T) {
			stubber := testtools.NewStubber()

			stubber.Add(makeGetParametersByPathStub([]types.Parameter{
				makeParam(&accessToken, baseParamPath+"/"+accessTokenPath),
				makeParam(&refreshToken, baseParamPath+"/"+refreshTokenPath),
			}))

			provider := setupProvider(t, stubber)

			token := &token{
				Value:     accessToken,
				Updatable: tt.accessTokenUpdated,
			}
			stubber.Add(testtools.Stub{
				OperationName: "PutParameter",
				Input: &ssm.PutParameterInput{
					Name:      utils.Ptr(baseParamPath + "/" + accessTokenPath),
					Value:     utils.Ptr(token.Value),
					Overwrite: utils.Ptr(true),
				},
				Output: &ssm.PutParameterOutput{},
			})

			err := provider.updateToken(token, token.Value, accessTokenPath)
			if tt.accessTokenUpdated {
				if err != nil {
					t.Errorf("updateToken() error = %v", err)
				}
			} else if err.Error() != tt.wantErr.Error() {
				t.Errorf("updateToken() error = %v, wantErr = %v", tt.wantErr, err)
			}
		})
	}
}

func setupProvider(t *testing.T, stubber *testtools.AwsmStubber) *TokensProvider {
	ps := ssm.NewFromConfig(*stubber.SdkConfig)
	pp := aws.NewSSMParamsProvider(ps, baseParamPath)
	if err := pp.PrefetchParams(); err != nil {
		t.Fatalf("couldn't init params provider: %s", err)
	}

	provider := NewTokensProvider(pp)
	return provider
}

func makeGetParametersByPathStub(params []types.Parameter) testtools.Stub {
	return testtools.Stub{
		OperationName: "GetParametersByPath",
		Input: &ssm.GetParametersByPathInput{
			Path:           utils.Ptr(baseParamPath),
			WithDecryption: utils.Ptr(true),
			Recursive:      utils.Ptr(true),
		},
		Output: &ssm.GetParametersByPathOutput{
			Parameters: params,
		},
	}
}

func makeParam(token *string, path string) types.Parameter {
	return types.Parameter{
		Name:  utils.Ptr(path),
		Value: token,
	}
}
