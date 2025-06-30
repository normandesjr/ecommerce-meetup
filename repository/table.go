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

func (d *dynamoDBRepo) CreateTable(ctx context.Context, action func()) error {
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
			{
				AttributeName: aws.String("GSI1PK"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("GSI1SK"),
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
		GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
			{
				IndexName: aws.String("GSI1"),
				OnDemandThroughput: &types.OnDemandThroughput{
					MaxReadRequestUnits:  aws.Int64(5),
					MaxWriteRequestUnits: aws.Int64(5),
				},
				Projection: &types.Projection{
					ProjectionType: types.ProjectionTypeAll,
				},
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("GSI1PK"),
						KeyType:       types.KeyTypeHash,
					},
					{
						AttributeName: aws.String("GSI1SK"),
						KeyType:       types.KeyTypeRange,
					},
				},
			},
		},
	}

	_, err := d.client.CreateTable(ctx, param)
	if err != nil {
		var resourceInUseErr *types.ResourceInUseException
		if errors.As(err, &resourceInUseErr) {
			return fmt.Errorf("%w", ErrTableAlreadyExists)
		}
		return err
	}

	waiterFunc := func(ctx context.Context) error {
		waiter := dynamodb.NewTableExistsWaiter(d.client)
		return waiter.Wait(ctx, &dynamodb.DescribeTableInput{
			TableName: aws.String(d.tableName)}, 5*time.Minute)
	}

	return d.waitForOperation(ctx, action, waiterFunc)
}

func (d *dynamoDBRepo) DeleteTable(ctx context.Context, action func()) error {
	_, err := d.client.DeleteTable(context.TODO(), &dynamodb.DeleteTableInput{
		TableName: aws.String(d.tableName),
	})
	if err != nil {
		return err
	}

	waiterFunc := func(ctx context.Context) error {
		waiter := dynamodb.NewTableNotExistsWaiter(d.client)
		return waiter.Wait(ctx, &dynamodb.DescribeTableInput{
			TableName: aws.String(d.tableName)}, 5*time.Minute)
	}

	return d.waitForOperation(ctx, action, waiterFunc)
}

func (d *dynamoDBRepo) waitForOperation(ctx context.Context, tickerAction func(), waiterFunc func(context.Context) error) error {
	errCh := make(chan error, 1)
	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()

	go func() {
		errCh <- waiterFunc(ctx)
	}()

	for {
		select {
		case err := <-errCh:
			return err
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			tickerAction()
		}
	}
}
