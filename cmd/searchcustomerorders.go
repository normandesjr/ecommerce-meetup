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

var searchCustomer string
var limit int

var searchCustomerOrdersCmd = &cobra.Command{
	Use:   "search-customer-orders",
	Short: "Search the last customer search orders",
	Run:   searchCustomerOrders,
}

func init() {
	searchCustomerOrdersCmd.Flags().StringVar(&searchCustomer, "customer", "", "The username to search")
	searchCustomerOrdersCmd.MarkFlagRequired("customer")

	searchCustomerOrdersCmd.Flags().IntVar(&limit, "limit", 2, "The limit to search orders")
	rootCmd.AddCommand(searchCustomerOrdersCmd)
}

func searchCustomerOrders(cmd *cobra.Command, args []string) {
	log.Println("Searching orders...")

	keyEx := expression.Key("PK").Equal(expression.Value(fmt.Sprintf("CUSTOMER#%s", searchCustomer)))
	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()
	if err != nil {
		log.Fatalf("Couldn't not build expression for query: %v", err)
	}
	param := &dynamodb.QueryInput{
		TableName:                 &Dynamo.TableName,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		Limit:                     aws.Int32(int32(limit) + 1),
		ScanIndexForward:          aws.Bool(false),
	}

	response, err := Dynamo.DynamoDbClient.Query(context.TODO(), param)
	if err != nil {
		log.Fatalf("Error searching customer orders: %v", err)
	}

	var orders []Order
	var customer Customer

	customerItem := response.Items[0]
	err = attributevalue.UnmarshalMap(customerItem, &customer)
	if err != nil {
		log.Fatalf("Error unmarshaling customer: %v", err)
	}

	err = attributevalue.UnmarshalListOfMaps(response.Items[1:], &orders)
	if err != nil {
		log.Fatalf("Error unmarshaling orders: %v", err)
	}

	fmt.Printf("\nCustomer [%s]\n", customer)

	fmt.Println("Orders:")
	for _, o := range orders {
		fmt.Println(o)
	}

}
