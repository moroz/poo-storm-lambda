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

func HandleCreateCommentRequest(d *dynamodb.Client, body string) (events.LambdaFunctionURLResponse, error) {
	var input models.CreateCommentInput
	err := json.Unmarshal([]byte(body), &input)
	if err != nil {
		return events.LambdaFunctionURLResponse{
			StatusCode: 400,
			Body:       `{"error":"Bad Request"}`,
		}, err
	}

	successResponse := events.LambdaFunctionURLResponse{
		StatusCode: 201,
		Body:       `{"status":"Created"}`,
	}

	// "I am a robot" is a CAPTCHA-like field in the
	// form, hidden to humans. If we are dealing with a robot,
	// we do not store comments and pretend that we did
	if input.IAmARobot {
		return successResponse, nil
	}

	err = models.CreateComment(d, &input)
	if err != nil {
		return events.LambdaFunctionURLResponse{
			StatusCode: 422,
			Body:       `{"error":"Unprocessable entity"}`,
		}, err
	}

	return successResponse, nil
}
