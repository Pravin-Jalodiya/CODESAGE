package db

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func DynamoInitClient() (*dynamodb.Client, error) {
	// Load the AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	// Create and return DynamoDB client
	return dynamodb.NewFromConfig(cfg), nil
}
