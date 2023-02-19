package api

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/moroz/poo-storm-lambda/models"
)

type ListCommentsRequestResponse struct {
	Comments *[]models.Comment `json:"comments"`
}

func HandleListCommentsRequest(d *dynamodb.Client, url string) (events.LambdaFunctionURLResponse, error) {
	comments, err := models.ListComments(d, url)
	if err != nil {
		return events.LambdaFunctionURLResponse{Body: "Failed to fetch comments", StatusCode: 500}, err
	}
	result := ListCommentsRequestResponse{
		Comments: comments,
	}
	json, err := json.Marshal(result)
	if err != nil {
		return events.LambdaFunctionURLResponse{Body: "Failed to serialize comments as JSON", StatusCode: 500}, err
	}

	return events.LambdaFunctionURLResponse{Body: string(json), StatusCode: 200}, nil
}
