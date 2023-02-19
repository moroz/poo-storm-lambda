package models

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/pkg/errors"
	uuid "github.com/uuid6/uuid6go-proto"
)

const TABLE_NAME = "Comments"

type CreateCommentInput struct {
	Body      string `json:"body"`
	Signature string `json:"signature"`
	Url       string `json:"url"`
}

type Comment struct {
	ID        uuid.UUIDv7 `dynamodbav:"id" json:"id"`
	Body      string      `dynamodbav:"body" json:"body"`
	Signature string      `dynamodbav:"signature" json:"signature"`
	Url       string      `dynamodbav:"url" json:"url"`
}

func (c Comment) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"id":        c.ID.ToString(),
		"body":      c.Body,
		"signature": c.Signature,
		"url":       c.Url,
	})
}

func CreateComment(d *dynamodb.Client, input *CreateCommentInput) error {
	generator := &uuid.UUIDv7Generator{}
	id := generator.Next()

	newComment := &Comment{
		ID:        id,
		Body:      input.Body,
		Signature: input.Signature,
		Url:       input.Url,
	}

	item, err := attributevalue.MarshalMap(newComment)
	if err != nil {
		log.Fatal(err)
	}

	request := &dynamodb.PutItemInput{
		TableName: aws.String(TABLE_NAME),
		Item:      item,
	}

	result, err := d.PutItem(context.TODO(), request)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(result)

	return nil
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
