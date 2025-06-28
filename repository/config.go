package repository

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type dynamoDBRepo struct {
	client    *dynamodb.Client
	tableName *string
}

func NewDynamoDBRepo(profile string) (*dynamoDBRepo, error) {
	config, err := config.LoadDefaultConfig(context.TODO(), func(lo *config.LoadOptions) error {
		lo.SharedConfigProfile = profile
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("connecting to AWS using the profile %q: %v", profile, err)
	}

	client := dynamodb.NewFromConfig(config)
	return &dynamoDBRepo{
		client:    client,
		tableName: aws.String("ECommerceCloudDay"),
	}, nil
}
