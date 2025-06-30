package cmd

import (
	"context"
	"fmt"
	"meetup/repository"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var deleteTableCmd = &cobra.Command{
	Use:          "delete-table",
	Aliases:      []string{"dt"},
	Short:        "Delete the DynamoDB table",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		profile := viper.GetString("profile")
		tableName := viper.GetString("table")

		return deleteTable(profile, tableName)
	},
}

func init() {
	rootCmd.AddCommand(deleteTableCmd)
}

func deleteTable(profile, tableName string) error {
	repo, err := repository.NewDynamoDBRepo(profile, tableName)
	if err != nil {
		return err
	}

	action := func() {
		fmt.Printf("...")
	}

	err = repo.DeleteTable(context.Background(), action)
	if err != nil {
		return err
	}

	fmt.Printf("\nTable %s deleted\n", tableName)

	return nil
}
