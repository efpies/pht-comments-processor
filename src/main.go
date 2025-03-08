package main

import (
	"encoding/json"
	awsLambda "github.com/aws/aws-lambda-go/lambda"
	"log"
	"pht/comments-processor/handlers/lambda"
	"pht/comments-processor/model"
	"strconv"
	"strings"
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
	le, err := tryParseLambdaEvent(event)
	if err != nil {
		return nil, err
	}

	if le != nil {
		log.Printf("Received event: %v -> %#v", string(event), le)
		return handleEvent(le)
	}

	re, err := tryParseRouteEvent(event)
	if err != nil {
		return nil, err
	}

	if re != nil {
		log.Printf("Received api call: %v -> %#v", string(event), re)
		return handleApi(re, err)
	}

	return nil, nil
}

func tryParseLambdaEvent(event json.RawMessage) (*lambda.Event, error) {
	var le *lambda.Event
	if err := json.Unmarshal(event, &le); err != nil {
		log.Fatalf("Failed to unmarshal event: %v", err)
		return nil, err
	}

	if le.Platform == "" {
		return nil, nil
	}
	return le, nil
}

func tryParseRouteEvent(event json.RawMessage) (*lambda.RouteEvent, error) {
	var re *lambda.RouteEvent
	if err := json.Unmarshal(event, &re); err != nil {
		log.Fatalf("Failed to unmarshal event: %v", err)
		return nil, err
	}

	if re.RawPath == "" {
		return nil, nil
	}

	return re, nil
}

func handleEvent(le *lambda.Event) (any, error) {
	result, err := appSvc.lambdaEventHandler.Handle(le)
	if err != nil {
		result = model.ErrorDto{Message: err.Error()}
	}

	return result, nil
}

func handleApi(re *lambda.RouteEvent, err error) (any, error) {
	flowID := strings.TrimPrefix(re.RawPath, "/")
	_, err = strconv.Atoi(flowID)
	if err != nil {
		return nil, nil
	}

	ndg := appSvc.phtServices.GetNotifierDataGetter()
	result, err := ndg.GetAllNotifierData(flowID)
	if err != nil {
		return nil, err
	}

	return map[string]any{
		"statusCode": 200,
		"headers": map[string]string{
			"Access-Control-Allow-Headers": "Content-Type",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "GET",
		},
		"body": result,
	}, nil
}
