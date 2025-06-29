package cmd

import (
	"errors"
	"fmt"
	"meetup/repository"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var createCustomerCmd = &cobra.Command{
	Use:          "create-customer username=<username> email=<email> name=<name>",
	Aliases:      []string{"cc"},
	Short:        "Save new customer to DynamoDB table",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		profile := viper.GetString("profile")
		tableName := viper.GetString("table")

		username, err := cmd.Flags().GetString("username")
		if err != nil {
			return err
		}
		email, err := cmd.Flags().GetString("email")
		if err != nil {
			return err
		}
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			return err
		}

		return createCustomer(profile, tableName, repository.Customer{
			Username: username,
			Email:    email,
			Name:     name,
		})
	},
}

func init() {
	createCustomerCmd.Flags().StringP("username", "u", "", "The username to save")
	createCustomerCmd.MarkFlagRequired("username")

	createCustomerCmd.Flags().StringP("email", "e", "", "The email to save")
	createCustomerCmd.MarkFlagRequired("email")

	createCustomerCmd.Flags().StringP("name", "n", "", "The name to save")
	createCustomerCmd.MarkFlagRequired("name")

	rootCmd.AddCommand(createCustomerCmd)
}

func createCustomer(profile, tableName string, customer repository.Customer) error {
	repo, err := repository.NewDynamoDBRepo(profile, tableName)
	if err != nil {
		return err
	}

	err = repo.CreateCustomer(customer)
	if errors.Is(err, repository.ErrCustomerAlreadyExists) {
		fmt.Println(err)
		return nil
	}

	return err
}
