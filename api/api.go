package api

import (
	"encoding/json"
	"log"

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

type CreateCommentResponse struct {
	Success bool            `json:"success"`
	Errors  []string        `json:"errors"`
	Comment *models.Comment `json:"comment"`
}

func createComment(d *dynamodb.Client, params *models.CreateCommentInput) (*CreateCommentResponse, error) {
	if params.IAmARobot {
		return &CreateCommentResponse{
			Success: true,
		}, nil
	}

	errors := params.Validate()
	if len(errors) != 0 {
		return &CreateCommentResponse{
			Success: false,
			Errors:  errors,
		}, nil
	}

	comment, err := models.CreateComment(d, params)
	if err != nil {
		return nil, err
	}

	return &CreateCommentResponse{
		Success: true,
		Comment: comment,
	}, nil
}

func invalidateCache(url string) {
	client := models.CreateCFClient()
	models.InvalidateCommentCache(client, url)
}

func HandleCreateCommentRequest(d *dynamodb.Client, body string) (events.LambdaFunctionURLResponse, error) {
	var input models.CreateCommentInput
	err := json.Unmarshal([]byte(body), &input)
	if err != nil {
		return events.LambdaFunctionURLResponse{
			StatusCode: 400,
			Body:       `{"success":false,"errors":["Bad Request"]}`,
		}, err
	}

	result, err := createComment(d, &input)
	if err != nil {
		log.Println(err)
		return events.LambdaFunctionURLResponse{
			StatusCode: 500,
			Body:       err.Error(),
		}, err
	}

	encoded, _ := json.Marshal(result)

	if result.Comment != nil {
		invalidateCache(input.Url)
	}

	status := 200
	if !result.Success {
		status = 422
	}

	return events.LambdaFunctionURLResponse{
		StatusCode: status,
		Body:       string(encoded),
	}, nil
}
