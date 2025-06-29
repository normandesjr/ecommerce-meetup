package cmd

import (
	"context"
	"errors"
	"fmt"
	"meetup/repository"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(createTableCmd)
}

var createTableCmd = &cobra.Command{
	Use:          "create-table",
	Aliases:      []string{"ct"},
	Short:        "Create the DynamoDB table",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		profile := viper.GetString("profile")
		tableName := viper.GetString("table")
		return createTable(profile, tableName)
	},
}

func createTable(profile, tableName string) error {
	repo, err := repository.NewDynamoDBRepo(profile, tableName)
	if err != nil {
		return err
	}

	err = repo.CreateTable(context.Background())
	if errors.Is(err, repository.ErrTableAlreadyExists) {
		fmt.Printf("Table %s already exists\n", tableName)
		return nil
	}

	fmt.Printf("Table %s created\n", tableName)

	return err
}
