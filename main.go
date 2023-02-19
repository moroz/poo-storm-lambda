package main

import (
	"fmt"
	"log"

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

func main() {
	fmt.Println("Hello world!")

	sample := &CreateCommentInput{
		Body:      "test",
		Signature: "km",
		Url:       "/blog/test",
	}

	err := CreateComment(client, sample)

	if err != nil {
		log.Fatal(err)
	}
}
