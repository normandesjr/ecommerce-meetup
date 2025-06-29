package cmd

import (
	"context"
	"meetup/repository"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var updateAddressCmd = &cobra.Command{
	Use:          "update-address --username <username> --id <address id> --street-address <street address> --zip-code <zip code>",
	Aliases:      []string{"ua"},
	Short:        "Update address to the customer",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		profile := viper.GetString("profile")
		tableName := viper.GetString("table")

		username, err := cmd.Flags().GetString("username")
		if err != nil {
			return err
		}
		addressId, err := cmd.Flags().GetString("id")
		if err != nil {
			return err
		}
		streetAddress, err := cmd.Flags().GetString("street-address")
		if err != nil {
			return err
		}
		zipCode, err := cmd.Flags().GetString("zip-code")
		if err != nil {
			return err
		}

		customer := repository.Customer{Username: username}
		address := repository.Address{
			Id:            addressId,
			StreetAddress: streetAddress, ZipCode: zipCode,
		}

		return addAddress(profile, tableName, customer, address)
	},
}

func init() {
	updateAddressCmd.Flags().StringP("username", "u", "", "The username to add the address")
	updateAddressCmd.MarkFlagRequired("username")

	updateAddressCmd.Flags().String("id", "", "The address identifier")
	updateAddressCmd.MarkFlagRequired("id")

	updateAddressCmd.Flags().StringP("street-address", "s", "", "The street address to save")
	updateAddressCmd.MarkFlagRequired("street-address")

	updateAddressCmd.Flags().StringP("zip-code", "z", "", "The zip code to save")
	updateAddressCmd.MarkFlagRequired("zip-code")

	rootCmd.AddCommand(updateAddressCmd)
}

func addAddress(profile, tableName string, customer repository.Customer, address repository.Address) error {
	repo, err := repository.NewDynamoDBRepo(profile, tableName)
	if err != nil {
		return err
	}

	return repo.UpdateAddress(context.Background(), customer, address)
}
