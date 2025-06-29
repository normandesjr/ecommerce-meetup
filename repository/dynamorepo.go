package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var (
	ErrTableAlreadyExists    = errors.New("table already exists")
	ErrCustomerAlreadyExists = errors.New("customer username or email already exists")
	ErrCustomerNotFound      = errors.New("customer not found")
)

type dynamoDBRepo struct {
	client    *dynamodb.Client
	tableName string
}

func NewDynamoDBRepo(profile, tableName string) (*dynamoDBRepo, error) {
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
		tableName: tableName,
	}, nil
}
