package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	uuid "github.com/uuid6/uuid6go-proto"
)

type CreateCommentInput struct {
	Body      string `json:"body"`
	Signature string `json:"signature"`
	Url       string `json:"url"`
}

type NewComment struct {
	ID        []byte `dynamodbav:"id"`
	Body      string `dynamodbav:"body"`
	Signature string `dynamodbav:"signature"`
	Url       string `dynamodbav:"url"`
}

func UUIDToBytes(uuid uuid.UUIDv7) []byte {
	result := make([]byte, 16)
	for i, b := range uuid {
		result[i] = b
	}
	return result
}

func CreateComment(d *dynamodb.Client, input *CreateCommentInput) error {
	generator := &uuid.UUIDv7Generator{}
	id := generator.Next()

	newComment := &NewComment{
		ID:        id[:],
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
