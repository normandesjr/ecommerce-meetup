package cmd

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(createTableCmd)
}

var createTableCmd = &cobra.Command{
	Use:   "create-table",
	Short: "Create the DynamoDB table",
	Run:   createTable,
}

func createTable(cmd *cobra.Command, args []string) {
	param := &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("PK"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("SK"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("PK"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("SK"),
				KeyType:       types.KeyTypeRange,
			},
		},
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
		TableName: aws.String(Dynamo.TableName),
	}

	_, err := Dynamo.DynamoDbClient.CreateTable(context.TODO(), param)
	if err != nil {
		log.Fatalf("Got error calling CreateTable: %v", err)
	}
	log.Println("Table is created")
}
