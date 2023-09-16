package cmd

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "meetup",
	Short: "Meetup AWS User Grupo Triangulo Mineiro",
}

func init() {
	cobra.OnInitialize(loadAWSConfig)
}

type tableBasics struct {
	DynamoDbClient *dynamodb.Client
	TableName      string
}

var Dynamo *tableBasics

func loadAWSConfig() {
	log.Println("loading aws config...")
	defaultConfig, err := config.LoadDefaultConfig(context.TODO(), func(lo *config.LoadOptions) error {
		lo.SharedConfigProfile = "soudev"
		return nil
	})
	if err != nil {
		log.Fatalf("error connecting to AWS %v", err)
	}

	client := dynamodb.NewFromConfig(defaultConfig)
	tableName := "EcommerceMeetup"

	Dynamo = &tableBasics{
		DynamoDbClient: client,
		TableName:      tableName,
	}

}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
