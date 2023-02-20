package main

import (
	"context"
	"log"
	"regexp"

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

var pathRegexp = regexp.MustCompile("^/comments/")

func extractUrlFromRawPath(path string) string {
	if !pathRegexp.Match([]byte(path)) {
		return ""
	}
	return "/" + pathRegexp.ReplaceAllString(path, "")
}

// Naive router
func handleRequest(ctx context.Context, request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	// non-null body and path == '/comments': treat as POST request
	if request.Body != "" && request.RawPath == "/comments" {
		return api.HandleCreateCommentRequest(client, request.Body)
	}
	url := extractUrlFromRawPath(request.RawPath)
	if url == "" {
		return events.LambdaFunctionURLResponse{
			StatusCode: 400,
			Body:       `{"status":"Bad Request"}`,
		}, nil
	}
	return api.HandleListCommentsRequest(client, url)
}

func main() {
	runtime.Start(handleRequest)
}
