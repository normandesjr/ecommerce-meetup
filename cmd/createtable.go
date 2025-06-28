package cmd

import (
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
		return createTable(profile)
	},
}

func createTable(profile string) error {
	repo, err := repository.NewDynamoDBRepo(profile)
	if err != nil {
		return err
	}

	return repo.CreateTable()
}
