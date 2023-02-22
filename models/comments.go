package models

import (
	"context"
	"encoding/json"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/pkg/errors"
	uuid "github.com/uuid6/uuid6go-proto"
)

const TABLE_NAME = "Comments"

type CreateCommentInput struct {
	Body      string  `json:"body"`
	Signature string  `json:"signature"`
	Url       string  `json:"url"`
	Website   *string `json:"website"`
	Email     string  `json:"email"`
	IAmARobot bool    `json:"iAmARobot"`
}

type Comment struct {
	ID         uuid.UUIDv7 `dynamodbav:"id" json:"id"`
	Email      string      `dynamodbav:"email" json:"-"`
	Body       string      `dynamodbav:"body" json:"body"`
	Signature  string      `dynamodbav:"signature" json:"signature"`
	Url        string      `dynamodbav:"url" json:"url"`
	Website    *string     `dynamodbav:"website" json:"website"`
	InsertedAt int64       `dynamodbav:"insertedAt" json:"insertedAt"`
}

func (c Comment) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"id":         c.ID.ToString(),
		"body":       c.Body,
		"signature":  c.Signature,
		"url":        c.Url,
		"website":    c.Website,
		"insertedAt": c.InsertedAt,
	})
}

func CreateComment(d *dynamodb.Client, input *CreateCommentInput) (*Comment, error) {
	generator := &uuid.UUIDv7Generator{}
	id := generator.Next()

	newComment := &Comment{
		ID:         id,
		Body:       input.Body,
		Email:      input.Email,
		Signature:  input.Signature,
		Url:        input.Url,
		Website:    input.Website,
		InsertedAt: time.Now().Unix(),
	}

	item, err := attributevalue.MarshalMap(newComment)
	if err != nil {
		return nil, err
	}

	request := &dynamodb.PutItemInput{
		TableName: aws.String(TABLE_NAME),
		Item:      item,
	}

	_, err = d.PutItem(context.TODO(), request)
	if err != nil {
		return nil, err
	}

	return newComment, nil
}

func ListComments(d *dynamodb.Client, url string) (*[]Comment, error) {
	keyEx := expression.Key("url").Equal(expression.Value(url))
	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to serialize key condition expression")
	}

	query := &dynamodb.QueryInput{
		TableName:                 aws.String(TABLE_NAME),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	}

	result, err := d.Query(context.TODO(), query)
	if err != nil {
		return nil, err
	}

	var comments []Comment
	attributevalue.UnmarshalListOfMaps(result.Items, &comments)

	return &comments, nil
}
