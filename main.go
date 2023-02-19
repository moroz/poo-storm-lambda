package main

import (
	"context"
	"fmt"
	"log"

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

// Naive router
func handleRequest(ctx context.Context, request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	// non-null body and path == '/comments': treat as POST request
	if request.Body != "" && request.RawPath == "/comments" {
		return events.LambdaFunctionURLResponse{Body: request.Body, StatusCode: 200}, nil
	} else {
		return HandleListCommentsRequest(client, request.QueryStringParameters["url"])
		// return events.LambdaFunctionURLResponse{Body: request.RawQueryString, StatusCode: 200}, nil
	}
}

// func main() {
// 	runtime.Start(handleRequest)
// }

func main() {
	res, err := ListComments(client, "/blog/test")
	if err != nil {
		log.Fatal(err)
	}

	for _, comment := range *res {
		fmt.Println(comment.Url, comment.Body, comment.ID)
	}
}
