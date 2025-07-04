package repository

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func (d *dynamoDBRepo) UpdateAddress(ctx context.Context, customer Customer, address Address) error {
	addressId := fmt.Sprintf("addresses.%s", address.Id)
	update := expression.Set(expression.Name(addressId), expression.Value(address))

	return d.updateItem(ctx, customer, update)
}

func (d *dynamoDBRepo) RemoveAddress(ctx context.Context, customer Customer, addressId string) error {
	id := fmt.Sprintf("addresses.%s", addressId)
	update := expression.Remove(expression.Name(id))

	return d.updateItem(ctx, customer, update)
}

func (d *dynamoDBRepo) updateItem(ctx context.Context, customer Customer, update expression.UpdateBuilder) error {
	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		return err
	}

	_, err = d.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:                 aws.String(d.tableName),
		Key:                       customer.GetKey(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
	})

	return err
}
