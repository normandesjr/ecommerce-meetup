package cmd

// import (
// 	"context"
// 	"fmt"
// 	"log"

// 	"github.com/aws/aws-sdk-go-v2/aws"
// 	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
// 	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
// 	"github.com/spf13/cobra"
// )

// var username string
// var email string
// var name string

// var createCustomerCmd = &cobra.Command{
// 	Use:   "create-customer",
// 	Short: "Save new Customer to DynamoDB",
// 	Run:   createCustomer,
// }

// func init() {
// 	createCustomerCmd.Flags().StringVar(&username, "username", "", "The username to save")
// 	createCustomerCmd.MarkFlagRequired("username")

// 	createCustomerCmd.Flags().StringVar(&email, "email", "", "The email to save")
// 	createCustomerCmd.MarkFlagRequired("email")

// 	createCustomerCmd.Flags().StringVar(&name, "name", "", "The name to save")
// 	createCustomerCmd.MarkFlagRequired("name")
// 	rootCmd.AddCommand(createCustomerCmd)
// }

// func createCustomer(cmd *cobra.Command, args []string) {
// 	log.Println("Saving customer...")
// 	param := &dynamodb.TransactWriteItemsInput{
// 		TransactItems: []types.TransactWriteItem{
// 			{
// 				Put: &types.Put{
// 					Item: map[string]types.AttributeValue{
// 						"PK":        &types.AttributeValueMemberS{Value: fmt.Sprintf("CUSTOMER#%s", username)},
// 						"SK":        &types.AttributeValueMemberS{Value: fmt.Sprintf("CUSTOMER#%s", username)},
// 						"Username":  &types.AttributeValueMemberS{Value: username},
// 						"Email":     &types.AttributeValueMemberS{Value: email},
// 						"Name":      &types.AttributeValueMemberS{Value: name},
// 						"Addresses": &types.AttributeValueMemberM{Value: nil},
// 					},
// 					TableName:           &Dynamo.TableName,
// 					ConditionExpression: aws.String("attribute_not_exists(PK)"),
// 				},
// 			},
// 			{
// 				Put: &types.Put{
// 					Item: map[string]types.AttributeValue{
// 						"PK":       &types.AttributeValueMemberS{Value: fmt.Sprintf("CUSTOMEREMAIL#%s", email)},
// 						"SK":       &types.AttributeValueMemberS{Value: fmt.Sprintf("CUSTOMEREMAIL#%s", email)},
// 						"Username": &types.AttributeValueMemberS{Value: username},
// 					},
// 					TableName:           &Dynamo.TableName,
// 					ConditionExpression: aws.String("attribute_not_exists(PK)"),
// 				},
// 			},
// 		},
// 	}

// 	_, err := Dynamo.DynamoDbClient.TransactWriteItems(context.TODO(), param)
// 	if err != nil {
// 		log.Fatalf("Error creating customer: %v", err)
// 	}

// 	log.Println("Customer created")

// }
