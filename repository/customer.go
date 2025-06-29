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
	customer.PK = fmt.Sprintf("CUSTOMER#%s", customer.Username)
	customer.SK = fmt.Sprintf("CUSTOMER#%s", customer.Username)
	customer.Addresses = make(map[string]Address)
	mashaledCustomer, err := attributevalue.MarshalMap(customer)
	if err != nil {
		return err
	}

	customerEmail := CustomerEmail{
		PK:       fmt.Sprintf("CUSTOMEREMAIL#%s", customer.Email),
		SK:       fmt.Sprintf("CUSTOMEREMAIL#%s", customer.Email),
		Username: customer.Username,
	}
	marshaledCustomerEmail, err := attributevalue.MarshalMap(customerEmail)
	if err != nil {
		return err
	}

	transactWriteItems := &dynamodb.TransactWriteItemsInput{
		TransactItems: []types.TransactWriteItem{
			{
				Put: &types.Put{
					Item:                mashaledCustomer,
					TableName:           aws.String(d.tableName),
					ConditionExpression: aws.String("attribute_not_exists(PK)"),
				},
			},
			{
				Put: &types.Put{
					Item:                marshaledCustomerEmail,
					TableName:           aws.String(d.tableName),
					ConditionExpression: aws.String("attribute_not_exists(PK)"),
				},
			},
		},
	}

	_, err = d.client.TransactWriteItems(ctx, transactWriteItems)
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
