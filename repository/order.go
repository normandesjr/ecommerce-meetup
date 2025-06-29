package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
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

	order := Order{
		PK:        fmt.Sprintf("CUSTOMER#%s", customer.Username),
		SK:        fmt.Sprintf("#ORDER#%s", orderId),
		Id:        orderId,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
		Status:    "PENDING",
		ShippedTo: addressToShip,
		Total:     total,
		GSI1PK:    fmt.Sprintf("ORDER#%s", orderId),
		GSI1SK:    fmt.Sprintf("ORDER#%s", orderId),
	}
	marshaledOrder, err := attributevalue.MarshalMap(order)
	if err != nil {
		return err
	}

	transactItems[0] = types.TransactWriteItem{
		Put: &types.Put{
			Item:                marshaledOrder,
			TableName:           aws.String(d.tableName),
			ConditionExpression: aws.String("attribute_not_exists(PK)"),
		},
	}

	for i, v := range items {
		item := OrderItem{
			PK:          fmt.Sprintf("ORDER#%s#ITEM#%s", orderId, v.Id),
			SK:          fmt.Sprintf("ORDER#%s#ITEM#%s", orderId, v.Id),
			OrderId:     orderId,
			Description: v.Description,
			Price:       v.Price,
			GSI1PK:      fmt.Sprintf("ORDER#%s", orderId),
			GSI1SK:      fmt.Sprintf("ITEM#%s", v.Id),
		}
		marshaledItem, err := attributevalue.MarshalMap(item)
		if err != nil {
			return err
		}

		transactItems[i+1] = types.TransactWriteItem{
			Put: &types.Put{
				Item:                marshaledItem,
				TableName:           aws.String(d.tableName),
				ConditionExpression: aws.String("attribute_not_exists(PK)"),
			},
		}
	}

	_, err = d.client.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{TransactItems: transactItems})
	if err != nil {
		return err
	}

	return nil
}

func (d *dynamoDBRepo) GetOrders(ctx context.Context, customer *Customer) ([]Order, error) {
	keyEx := expression.Key("PK").Equal(expression.Value(fmt.Sprintf("CUSTOMER#%s", customer.Username))).
		And(expression.Key("SK").BeginsWith("#ORDER#"))
	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()
	if err != nil {
		return nil, err
	}

	queryPaginator := dynamodb.NewQueryPaginator(d.client, &dynamodb.QueryInput{
		TableName:                 aws.String(d.tableName),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	})

	var orders []Order
	for queryPaginator.HasMorePages() {
		response, err := queryPaginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		var orderPage []Order
		err = attributevalue.UnmarshalListOfMaps(response.Items, &orderPage)
		if err != nil {
			return nil, err
		}

		orders = append(orders, orderPage...)
	}

	return orders, nil
}
