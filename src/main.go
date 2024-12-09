package main

import (
	"encoding/json"
	awsLambda "github.com/aws/aws-lambda-go/lambda"
	"log"
	"pht/comments-processor/handlers/lambda"
	"pht/comments-processor/model"
)

var appSvc *appServices

func init() {
	var err error
	appSvc, err = newAppServices()
	if err != nil {
		panic(err)
	}

	if err = appSvc.init(); err != nil {
		panic(err)
	}
}

func main() {
	awsLambda.Start(handleRequest)
}

func handleRequest(event json.RawMessage) (any, error) {
	var le *lambda.Event
	if err := json.Unmarshal(event, &le); err != nil {
		log.Fatalf("Failed to unmarshal event: %v", err)
		return nil, err
	}

	log.Printf("Received event: %v -> %#v", string(event), le)

	return handleEvent(le)
}

func handleEvent(le *lambda.Event) (any, error) {
	result, err := appSvc.lambdaEventHandler.Handle(le)
	if err != nil {
		result = model.ErrorDto{Message: err.Error()}
	}

	return result, nil
}
