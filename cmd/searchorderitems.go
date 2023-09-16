package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/spf13/cobra"
)

var orderId string

var searchOrderItemsCmd = &cobra.Command{
	Use:   "search-order-items",
	Short: "Search the order items",
	Run:   searchOrderItems,
}

func init() {
	searchOrderItemsCmd.Flags().StringVar(&orderId, "order-id", "", "The orderId to search")
	searchOrderItemsCmd.MarkFlagRequired("order-id")

	rootCmd.AddCommand(searchOrderItemsCmd)
}

func searchOrderItems(cmd *cobra.Command, args []string) {
	log.Println("Searching order items...")

	keyEx := expression.Key("GSI1PK").Equal(expression.Value(fmt.Sprintf("ORDER#%s", orderId)))
	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()
	if err != nil {
		log.Fatalf("Couldn't not build expression for query: %v", err)
	}
	param := &dynamodb.QueryInput{
		TableName:                 &Dynamo.TableName,
		IndexName:                 aws.String("GSI1"),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		ScanIndexForward:          aws.Bool(false),
	}

	response, err := Dynamo.DynamoDbClient.Query(context.TODO(), param)
	if err != nil {
		log.Fatalf("Error searching customer orders: %v", err)
	}

	var order Order
	var orderItens []OrderItem

	err = attributevalue.UnmarshalMap(response.Items[0], &order)
	if err != nil {
		log.Fatalf("Error unmarshaling order: %v", err)
	}

	err = attributevalue.UnmarshalListOfMaps(response.Items[1:], &orderItens)
	if err != nil {
		log.Fatalf("Error unmarshaling orderItems: %v", err)
	}

	fmt.Printf("\nOrder [%s]\n", order)

	fmt.Println("Items:")
	for _, o := range orderItens {
		fmt.Println(o)
	}

}
