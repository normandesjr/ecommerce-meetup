package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var createOrderCmd = &cobra.Command{
	Use:          "create-order",
	Aliases:      []string{"co"},
	Short:        "Create an order for the choosed items",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		profile := viper.GetString("profile")
		tableName := viper.GetString("table")

		username, err := cmd.Flags().GetString("username")
		if err != nil {
			return err
		}
		shipAddress, err := cmd.Flags().GetString("ship-address")
		if err != nil {
			return err
		}
		items, err := cmd.Flags().GetStringSlice("items")
		if err != nil {
			return err
		}

		fmt.Println(username)
		fmt.Println(shipAddress)
		fmt.Printf("%d: %v\n", len(items), items)

		return createOrder(profile, tableName)
	},
}

func init() {
	createOrderCmd.Flags().StringP("username", "u", "", "The username to save")
	createOrderCmd.MarkFlagRequired("username")

	createOrderCmd.Flags().StringP("ship-address", "a", "", "The address id to ship the order")
	createOrderCmd.MarkFlagRequired("ship-address")

	createOrderCmd.Flags().StringSliceP("items", "i", nil, "The items id as comma list to add")
	createOrderCmd.MarkFlagRequired("items")

	rootCmd.AddCommand(createOrderCmd)
}

func createOrder(profile, tableName string) error {

	return nil
}

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
