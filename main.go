package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	runtime "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/moroz/poo-storm-lambda/api"
	"github.com/moroz/poo-storm-lambda/models"
)

var tableCreated bool = false

var client *dynamodb.Client

func init() {
	if tableCreated {
		return
	}
	log.Println("Creating table...")
	client = models.CreateClient()
	models.CreateTableIfNotExists(client)
	tableCreated = true
}

// Naive router
func handleRequest(ctx context.Context, request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	// non-null body and path == '/comments': treat as POST request
	if request.Body != "" && request.RawPath == "/comments" {
		return api.HandleCreateCommentRequest(client, request.Body)
	} else {
		return api.HandleListCommentsRequest(client, request.QueryStringParameters["url"])
	}
}

func main() {
	runtime.Start(handleRequest)
}
