package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func (d *dynamoDBRepo) CreateCustomer(ctx context.Context, customer Customer) error {
	transactWriteItems := &dynamodb.TransactWriteItemsInput{
		TransactItems: []types.TransactWriteItem{
			{
				// TODO: Deve ter como usar o dynamodbav da entidade... melhorar.
				Put: &types.Put{
					Item: map[string]types.AttributeValue{
						"PK":        &types.AttributeValueMemberS{Value: fmt.Sprintf("CUSTOMER#%s", customer.Username)},
						"SK":        &types.AttributeValueMemberS{Value: fmt.Sprintf("CUSTOMER#%s", customer.Username)},
						"username":  &types.AttributeValueMemberS{Value: customer.Username},
						"email":     &types.AttributeValueMemberS{Value: customer.Email},
						"name":      &types.AttributeValueMemberS{Value: customer.Name},
						"addresses": &types.AttributeValueMemberM{Value: nil},
					},
					TableName:           aws.String(d.tableName),
					ConditionExpression: aws.String("attribute_not_exists(PK)"),
				},
			},
			{
				Put: &types.Put{
					Item: map[string]types.AttributeValue{
						"PK":       &types.AttributeValueMemberS{Value: fmt.Sprintf("CUSTOMEREMAIL#%s", customer.Email)},
						"SK":       &types.AttributeValueMemberS{Value: fmt.Sprintf("CUSTOMEREMAIL#%s", customer.Email)},
						"username": &types.AttributeValueMemberS{Value: customer.Username},
					},
					TableName:           aws.String(d.tableName),
					ConditionExpression: aws.String("attribute_not_exists(PK)"),
				},
			},
		},
	}

	_, err := d.client.TransactWriteItems(ctx, transactWriteItems)
	if err != nil {
		var transactionCanceledException *types.TransactionCanceledException
		if errors.As(err, &transactionCanceledException) {
			cancellationReasons := transactionCanceledException.CancellationReasons
			for _, c := range cancellationReasons {
				if *c.Code == "ConditionalCheckFailed" {
					return fmt.Errorf("%w", ErrCustomerAlreadyExists)
				}
			}
		}

		return err
	}

	return nil
}

func (d *dynamoDBRepo) GetCustomer(ctx context.Context, username string) (*Customer, error) {
	customer := Customer{Username: username}
	response, err := d.client.GetItem(ctx, &dynamodb.GetItemInput{
		Key:                  customer.GetKey(),
		ProjectionExpression: aws.String("username"),
		TableName:            aws.String(d.tableName),
	})
	if err != nil {
		return nil, err
	}

	if len(response.Item) == 0 {
		return nil, ErrCustomerNotFound
	}

	err = attributevalue.UnmarshalMap(response.Item, &customer)
	if err != nil {
		return nil, err
	}

	return &customer, nil
}
