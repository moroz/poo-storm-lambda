package models

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	uuid "github.com/uuid6/uuid6go-proto"
)

var DISTRIBUTION_ID string

func init() {
	DISTRIBUTION_ID = os.Getenv("CLOUDFRONT_DISTRIBUTION_ID")
	if DISTRIBUTION_ID == "" {
		log.Fatalf("Environment variable CLOUDFRONT_DISTRIBUTION_ID not set!")
	}
}

func CreateCFClient() *cloudfront.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		panic(err)
	}

	return cloudfront.NewFromConfig(cfg)
}

func encodeCommentPath(path string) string {
	return "/comments" + path
}

func InvalidateCommentCache(cf *cloudfront.Client, path string) error {
	encoded := encodeCommentPath(path)

	uuidGenerator := uuid.UUIDv7Generator{}

	input := &cloudfront.CreateInvalidationInput{
		DistributionId: aws.String(DISTRIBUTION_ID),
		InvalidationBatch: &types.InvalidationBatch{
			CallerReference: aws.String(uuidGenerator.Next().ToString()),
			Paths: &types.Paths{
				Quantity: aws.Int32(1),
				Items:    aws.ToStringSlice([]*string{aws.String(encoded)}),
			},
		},
	}

	_, err := cf.CreateInvalidation(context.TODO(), input)
	if err != nil {
		return err
	}

	return nil
}
