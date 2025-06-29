package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func (d *dynamoDBRepo) CreateTable(ctx context.Context) error {
	param := &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("PK"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("SK"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("PK"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("SK"),
				KeyType:       types.KeyTypeRange,
			},
		},
		BillingMode: types.BillingModePayPerRequest,
		TableName:   aws.String(d.tableName),
	}

	_, err := d.client.CreateTable(ctx, param)
	if err != nil {
		var resourceInUseErr *types.ResourceInUseException
		if errors.As(err, &resourceInUseErr) {
			return fmt.Errorf("%w", ErrTableAlreadyExists)
		}
		return err
	}

	waiter := dynamodb.NewTableExistsWaiter(d.client)
	err = waiter.Wait(ctx, &dynamodb.DescribeTableInput{
		TableName: aws.String(d.tableName)}, 5*time.Minute)
	if err != nil {
		return err
	}

	return nil
}
