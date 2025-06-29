package repository

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func (d *dynamoDBRepo) GetOrderItems(ctx context.Context, orderId string) (OrderItems, error) {
	keyEx := expression.Key("GSI1PK").Equal(expression.Value(fmt.Sprintf("ORDER#%s", orderId))).
		And(expression.Key("GSI1SK").BeginsWith("ITEM#"))

	proj := expression.NamesList(
		expression.Name("description"),
		expression.Name("itemId"),
		expression.Name("price"),
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
		IndexName:                 aws.String("GSI1"),
	})

	var orderItems OrderItems
	for queryPaginator.HasMorePages() {
		response, err := queryPaginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		var orderItemsPage OrderItems
		err = attributevalue.UnmarshalListOfMaps(response.Items, &orderItemsPage)
		if err != nil {
			return nil, err
		}

		orderItems = append(orderItems, orderItemsPage...)
	}

	return orderItems, nil
}
