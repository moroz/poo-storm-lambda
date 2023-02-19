package main

import (
	"context"
	"log"

	runtime "github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var tableCreated bool = false

var client *dynamodb.Client

func init() {
	if tableCreated {
		return
	}
	log.Println("Creating table...")
	client = CreateClient()
	CreateTableIfNotExists(client)
	tableCreated = true
}

func handleRequest(ctx context.Context, request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	if request.Body != "" {
		return events.LambdaFunctionURLResponse{Body: request.Body, StatusCode: 200}, nil
	} else {
		return events.LambdaFunctionURLResponse{Body: request.RawQueryString, StatusCode: 200}, nil
	}
}

func main() {
	runtime.Start(handleRequest)
}
