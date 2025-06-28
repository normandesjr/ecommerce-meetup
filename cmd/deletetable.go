package cmd

// import (
// 	"context"
// 	"log"

// 	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
// 	"github.com/spf13/cobra"
// )

// func init() {
// 	rootCmd.AddCommand(deleteTableCmd)
// }

// var deleteTableCmd = &cobra.Command{
// 	Use:   "delete-table",
// 	Short: "Delete the DynamoDB table",
// 	Run:   deleteTable,
// }

// func deleteTable(cmd *cobra.Command, args []string) {
// 	_, err := Dynamo.DynamoDbClient.DeleteTable(context.TODO(), &dynamodb.DeleteTableInput{
// 		TableName: &Dynamo.TableName,
// 	})
// 	if err != nil {
// 		log.Fatalf("Got error calling DeleteTable: %v", err)
// 	}
// 	log.Println("Table is deleted")
// }
