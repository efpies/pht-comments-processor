package services

import (
	"pht/comments-processor/handlers/lambda"
	"testing"
)

const expectedResult = "mocked result"

func TestLambdaEventHandler_Handle(t *testing.T) {
	cases := []struct {
		platformRegistered bool
		wantResult         any
		wantErr            bool
	}{
		{
			platformRegistered: true,
			wantResult:         expectedResult,
			wantErr:            false,
		},
		{
			platformRegistered: false,
			wantResult:         nil,
			wantErr:            true,
		},
	}

	for _, tt := range cases {
		t.Run("Test LambdaEventHandler Handle", func(t *testing.T) {
			event := &lambda.Event{
				Platform:       "test",
				ServiceRequest: lambda.ServiceRequest{},
			}

			leh := NewLambdaEventHandler()
			if tt.platformRegistered {
				if err := leh.RegisterPlatformHandler("test", &serviceRequestHandlerMock{}); err != nil {
					t.Errorf("RegisterPlatformHandler() error = %v", err)
				}
			}

			result, err := leh.Handle(event)
			if tt.wantErr != (err != nil) {
				t.Errorf("Handle() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if tt.wantResult != result {
				t.Errorf("Handle() result = %v, want = %v", result, tt.wantResult)
			}
		})
	}
}

type serviceRequestHandlerMock struct {
}

func (m *serviceRequestHandlerMock) Handle(*lambda.ServiceRequest) (any, error) {
	return expectedResult, nil
}
