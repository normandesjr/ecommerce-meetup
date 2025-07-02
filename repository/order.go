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
		Username:  customer.Username,
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
			Id:          v.Id,
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

func (d *dynamoDBRepo) GetOrderById(ctx context.Context, orderId string) (*Order, error) {
	keyEx := expression.Key("GSI1PK").Equal(expression.Value(fmt.Sprintf("ORDER#%s", orderId))).
		And(expression.Key("GSI1SK").Equal(expression.Value(fmt.Sprintf("ORDER#%s", orderId))))

	expr, err := expression.NewBuilder().
		WithKeyCondition(keyEx).
		Build()
	if err != nil {
		return nil, err
	}

	response, err := d.client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(d.tableName),
		IndexName:                 aws.String("GSI1"),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	})
	if err != nil {
		return nil, err
	}

	if len(response.Items) == 0 {
		return nil, err
	}

	var foundOrder Order
	err = attributevalue.UnmarshalMap(response.Items[0], &foundOrder)
	if err != nil {
		return nil, err
	}

	return &foundOrder, nil

}

func (d *dynamoDBRepo) GetOrders(ctx context.Context, customer *Customer) ([]Order, error) {
	keyEx := expression.Key("PK").Equal(expression.Value(fmt.Sprintf("CUSTOMER#%s", customer.Username))).
		And(expression.Key("SK").BeginsWith("#ORDER#"))

	proj := expression.NamesList(
		expression.Name("createdAt"),
		expression.Name("orderId"),
		expression.Name("shippedTo"),
		expression.Name("status"),
		expression.Name("total"),
	)

	expr, err := expression.NewBuilder().
		WithKeyCondition(keyEx).
		WithProjection(proj).
		Build()
	if err != nil {
		return nil, err
	}

	queryPaginator := dynamodb.NewQueryPaginator(d.client, &dynamodb.QueryInput{
		TableName:                 aws.String(d.tableName),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		ProjectionExpression:      expr.Projection(),
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

func (d *dynamoDBRepo) UpdateStatusOrder(ctx context.Context, order *Order, status string) error {
	pk, err := attributevalue.Marshal(fmt.Sprintf("CUSTOMER#%s", order.Username))
	if err != nil {
		return err
	}
	sk, err := attributevalue.Marshal(fmt.Sprintf("#ORDER#%s", order.Id))
	if err != nil {
		return err
	}
	key := map[string]types.AttributeValue{"PK": pk, "SK": sk}

	update := expression.Set(expression.Name("status"), expression.Value(status))

	expr, err := expression.NewBuilder().
		WithUpdate(update).
		Build()
	if err != nil {
		return err
	}

	_, err = d.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:                 aws.String(d.tableName),
		Key:                       key,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
	})
	if err != nil {
		return err
	}

	return nil
}
