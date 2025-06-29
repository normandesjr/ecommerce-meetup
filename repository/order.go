package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/segmentio/ksuid"
)

func (d *dynamoDBRepo) CreateOrder(ctx context.Context, customer *Customer, items OrderItems) error {
	orderId := ksuid.New().String()

	var addressToShip string
	for a := range customer.Addresses {
		addressToShip = a
		break
	}

	total := items.Total()

	transactItems := make([]types.TransactWriteItem, len(items)+1)
	transactItems[0] = types.TransactWriteItem{
		Put: &types.Put{
			Item: map[string]types.AttributeValue{
				"PK":        &types.AttributeValueMemberS{Value: fmt.Sprintf("CUSTOMER#%s", customer.Username)},
				"SK":        &types.AttributeValueMemberS{Value: fmt.Sprintf("#ORDER#%s", orderId)},
				"orderId":   &types.AttributeValueMemberS{Value: orderId},
				"createdAt": &types.AttributeValueMemberS{Value: time.Now().UTC().Format(time.RFC3339)},
				"status":    &types.AttributeValueMemberS{Value: "PENDING"},
				"shippedTo": &types.AttributeValueMemberS{Value: addressToShip},
				"total":     &types.AttributeValueMemberN{Value: strconv.Itoa(total)},
				"GSI1PK":    &types.AttributeValueMemberS{Value: fmt.Sprintf("ORDER#%s", orderId)},
				"GSI1SK":    &types.AttributeValueMemberS{Value: fmt.Sprintf("ORDER#%s", orderId)},
			},
			TableName:           aws.String(d.tableName),
			ConditionExpression: aws.String("attribute_not_exists(PK)"),
		},
	}

	for i := 0; i < len(items); i++ {
		item := items[i]
		transactItems[i+1] = types.TransactWriteItem{
			Put: &types.Put{
				Item: map[string]types.AttributeValue{
					"PK":          &types.AttributeValueMemberS{Value: fmt.Sprintf("ORDER#%s#ITEM#%s", orderId, item.Id)},
					"SK":          &types.AttributeValueMemberS{Value: fmt.Sprintf("ORDER#%s#ITEM#%s", orderId, item.Id)},
					"orderId":     &types.AttributeValueMemberS{Value: orderId},
					"itemId":      &types.AttributeValueMemberS{Value: item.Id},
					"description": &types.AttributeValueMemberS{Value: item.Description},
					"price":       &types.AttributeValueMemberN{Value: strconv.Itoa(item.Price)},
					"GSI1PK":      &types.AttributeValueMemberS{Value: fmt.Sprintf("ORDER#%s", orderId)},
					"GSI1SK":      &types.AttributeValueMemberS{Value: fmt.Sprintf("ITEM#%s", item.Id)},
				},
				TableName:           aws.String(d.tableName),
				ConditionExpression: aws.String("attribute_not_exists(PK)"),
			},
		}
	}

	_, err := d.client.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{TransactItems: transactItems})
	if err != nil {
		return err
	}

	return nil
}
