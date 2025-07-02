package cmd

import (
	"context"
	"meetup/repository"

	"github.com/spf13/cobra"
)

func newRemoveAddressCmd(app *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "remove-address --username <username> --id <address id>",
		Aliases:      []string{"ra"},
		Short:        "Remove customer's address",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			username, _ := cmd.Flags().GetString("username")
			addressId, _ := cmd.Flags().GetString("id")

			customer := repository.Customer{Username: username}
			return removeAddress(app.repo, customer, addressId)
		},
	}

	cmd.Flags().StringP("username", "u", "", "The username to add the address")
	cmd.MarkFlagRequired("username")

	cmd.Flags().String("id", "", "The address identifier")
	cmd.MarkFlagRequired("id")

	return cmd
}

func removeAddress(repo Repository, customer repository.Customer, addressId string) error {
	return repo.RemoveAddress(context.Background(), customer, addressId)
}
