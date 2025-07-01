package cmd

import (
	"context"
	"errors"
	"fmt"
	"meetup/repository"

	"github.com/spf13/cobra"
)

func newCreateCustomerCmd(app *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "create-customer",
		Aliases:      []string{"cc"},
		Short:        "Save new customer to DynamoDB table",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			username, _ := cmd.Flags().GetString("username")
			email, _ := cmd.Flags().GetString("email")
			name, _ := cmd.Flags().GetString("name")

			return createCustomer(app.repo, repository.Customer{
				Username: username,
				Email:    email,
				Name:     name,
			})
		},
	}

	cmd.Flags().StringP("username", "u", "", "The username to save")
	cmd.MarkFlagRequired("username")

	cmd.Flags().StringP("email", "e", "", "The email to save")
	cmd.MarkFlagRequired("email")

	cmd.Flags().StringP("name", "n", "", "The name to save")
	cmd.MarkFlagRequired("name")

	return cmd
}

func createCustomer(repo Repository, customer repository.Customer) error {
	err := repo.CreateCustomer(context.Background(), customer)
	if errors.Is(err, repository.ErrCustomerAlreadyExists) {
		fmt.Println(err)
		return nil
	}

	return err
}
