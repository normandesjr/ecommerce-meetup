package cmd

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"math/rand"
// 	"time"

// 	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
// 	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
// 	"github.com/segmentio/ksuid"
// 	"github.com/spf13/cobra"
// )

// var customer string
// var amount float64
// var addItems bool

// func init() {
// 	createOrderCmd.Flags().StringVar(&customer, "customer", "", "The customer to save the order")
// 	createOrderCmd.MarkFlagRequired("customer")

// 	createOrderCmd.Flags().Float64Var(&amount, "amount", 0, "The order amount")
// 	createOrderCmd.MarkFlagRequired("amount")

// 	createOrderCmd.Flags().BoolVar(&addItems, "add-items", false, "Add 2 items to order")
// 	rootCmd.AddCommand(createOrderCmd)
// }

// var createOrderCmd = &cobra.Command{
// 	Use:   "create-order",
// 	Short: "Create an order",
// 	Run:   createOrder,
// }

// func createOrder(cmd *cobra.Command, args []string) {
// 	log.Println("Creating order...")

// 	orderId := ksuid.New().String()
// 	createdAt := time.Now().UTC().Format(time.RFC3339)

// 	param := &dynamodb.PutItemInput{
// 		Item: map[string]types.AttributeValue{
// 			"PK":        &types.AttributeValueMemberS{Value: fmt.Sprintf("CUSTOMER#%s", customer)},
// 			"SK":        &types.AttributeValueMemberS{Value: fmt.Sprintf("#ORDER#%s", orderId)},
// 			"OrderId":   &types.AttributeValueMemberS{Value: orderId},
// 			"CreatedAt": &types.AttributeValueMemberS{Value: createdAt},
// 			"Status":    &types.AttributeValueMemberS{Value: "PENDING"},
// 			"Amount":    &types.AttributeValueMemberN{Value: fmt.Sprintf("%.2f", amount)},
// 		},
// 		TableName: &Dynamo.TableName,
// 	}

// 	if addItems {
// 		addItemsToOrder(param, orderId)
// 	}

// 	_, err := Dynamo.DynamoDbClient.PutItem(context.TODO(), param)
// 	if err != nil {
// 		log.Fatalf("Error creating order: %v", err)
// 	}

// 	log.Println("Order created")
// }

// func addItemsToOrder(input *dynamodb.PutItemInput, orderId string) {
// 	input.Item["GSI1PK"] = &types.AttributeValueMemberS{Value: fmt.Sprintf("ORDER#%s", orderId)}
// 	input.Item["GSI1SK"] = &types.AttributeValueMemberS{Value: fmt.Sprintf("ORDER#%s", orderId)}

// 	// ATTENTION: I'M SAVING ITEMS OUTSIDE TRANSACTION TO BE ABLE DEMONSTRATE STEP BY STEP DURING THE TALK

// 	itens := []string{"Creatina", "Whey"}
// 	for i := 0; i < 2; i++ {
// 		itemId := fmt.Sprintf("%d", i+1)
// 		price := rand.Float64() * 100
// 		param := &dynamodb.PutItemInput{
// 			Item: map[string]types.AttributeValue{
// 				"PK":          &types.AttributeValueMemberS{Value: fmt.Sprintf("ORDER#%s#ITEM#%s", orderId, itemId)},
// 				"SK":          &types.AttributeValueMemberS{Value: fmt.Sprintf("ORDER#%s#ITEM#%s", orderId, itemId)},
// 				"OrderId":     &types.AttributeValueMemberS{Value: orderId},
// 				"ItemId":      &types.AttributeValueMemberS{Value: itemId},
// 				"Description": &types.AttributeValueMemberS{Value: itens[i]},
// 				"Price":       &types.AttributeValueMemberN{Value: fmt.Sprintf("%.2f", price)},
// 				"GSI1PK":      &types.AttributeValueMemberS{Value: fmt.Sprintf("ORDER#%s", orderId)},
// 				"GSI1SK":      &types.AttributeValueMemberS{Value: fmt.Sprintf("ITEM#%s", itemId)},
// 			},
// 			TableName: &Dynamo.TableName,
// 		}

// 		_, err := Dynamo.DynamoDbClient.PutItem(context.TODO(), param)
// 		if err != nil {
// 			log.Fatalf("Error creating order: %v", err)
// 		}
// 	}

// }
