package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const TABLE_NAME = "Comments"

func tableExists(d *dynamodb.Client, name string) bool {
	tables, err := d.ListTables(context.TODO(), &dynamodb.ListTablesInput{})
	if err != nil {
		log.Fatal("ListTables failed", err)
	}

	for _, n := range tables.TableNames {
		if n == name {
			return true
		}
	}

	return false
}

func buildCreateCommentsTableInput() *dynamodb.CreateTableInput {
	return &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("url"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("id"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("url"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("id"),
				KeyType:       types.KeyTypeRange,
			},
		},
		TableName:   aws.String(TABLE_NAME),
		BillingMode: types.BillingModePayPerRequest,
	}
}

func CreateTable(d *dynamodb.Client) error {
	_, err := d.CreateTable(context.TODO(), buildCreateCommentsTableInput())
	return err
}

func CreateTableIfNotExists(d *dynamodb.Client) error {
	if tableExists(d, TABLE_NAME) {
		return nil
	}
	return CreateTable(d)
}
