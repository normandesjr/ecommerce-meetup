package cmd

import (
	"context"
	"meetup/repository"

	"github.com/spf13/cobra"
)

func newUpdateAddressCmd(app *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "update-address --username <username> --id <address id> --street-address <street address> --zip-code <zip code>",
		Aliases:      []string{"ua"},
		Short:        "Update address to the customer",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			username, _ := cmd.Flags().GetString("username")
			addressId, _ := cmd.Flags().GetString("id")
			streetAddress, _ := cmd.Flags().GetString("street-address")
			zipCode, _ := cmd.Flags().GetString("zip-code")

			customer := repository.Customer{Username: username}
			address := repository.Address{
				Id:            addressId,
				StreetAddress: streetAddress, ZipCode: zipCode,
			}

			return addAddress(app.repo, customer, address)
		},
	}

	cmd.Flags().StringP("username", "u", "", "The username to add the address")
	cmd.MarkFlagRequired("username")

	cmd.Flags().String("id", "", "The address identifier")
	cmd.MarkFlagRequired("id")

	cmd.Flags().StringP("street-address", "s", "", "The street address to save")
	cmd.MarkFlagRequired("street-address")

	cmd.Flags().StringP("zip-code", "z", "", "The zip code to save")
	cmd.MarkFlagRequired("zip-code")

	return cmd
}

func addAddress(repo Repository, customer repository.Customer, address repository.Address) error {
	return repo.UpdateAddress(context.Background(), customer, address)
}
