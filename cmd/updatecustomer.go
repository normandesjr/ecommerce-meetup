package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/spf13/cobra"
)

var customerUpdate string
var addressId string
var address string

var updateCustomerCmd = &cobra.Command{
	Use:   "update-customer",
	Short: "Update customer address",
	Run:   updateCustomer,
}

func init() {
	updateCustomerCmd.Flags().StringVar(&customerUpdate, "customer", "", "The customer to update")
	updateCustomerCmd.MarkFlagRequired("customer")

	updateCustomerCmd.Flags().StringVar(&addressId, "addressId", "", "The address identifier")
	updateCustomerCmd.MarkFlagRequired("addressId")

	updateCustomerCmd.Flags().StringVar(&address, "address", "", "The address to save")
	updateCustomerCmd.MarkFlagRequired("address")

	rootCmd.AddCommand(updateCustomerCmd)
}

func updateCustomer(cmd *cobra.Command, args []string) {
	log.Println("Updating customer...")

	key, err := attributevalue.Marshal(fmt.Sprintf("CUSTOMER#%s", customerUpdate))
	if err != nil {
		log.Fatalf("Error creating key: %v", err)
	}

	addressKey := fmt.Sprintf("Addresses.%s", addressId)
	update := expression.Set(expression.Name(addressKey), expression.Value(address))
	condition := expression.AttributeExists(expression.Name("Addresses"))
	expr, err := expression.NewBuilder().WithCondition(condition).WithUpdate(update).Build()

	if err != nil {
		log.Fatalf("Error building the update expression: %v", err)
	}

	param := &dynamodb.UpdateItemInput{
		TableName:                 &Dynamo.TableName,
		Key:                       map[string]types.AttributeValue{"PK": key, "SK": key},
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
		ConditionExpression:       expr.Condition(),
	}

	_, err = Dynamo.DynamoDbClient.UpdateItem(context.TODO(), param)
	if err != nil {
		log.Fatalf("Error updating customer: %v", err)
	}

	log.Println("Customer updated")

}
