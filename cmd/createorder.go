package cmd

import (
	"context"
	"meetup/repository"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var createOrderCmd = &cobra.Command{
	Use:          "create-order",
	Aliases:      []string{"co"},
	Short:        "Create an order for the choosed items",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		profile := viper.GetString("profile")
		tableName := viper.GetString("table")

		username, err := cmd.Flags().GetString("username")
		if err != nil {
			return err
		}
		shipAddress, err := cmd.Flags().GetString("ship-address")
		if err != nil {
			return err
		}
		itemsId, err := cmd.Flags().GetIntSlice("items")
		if err != nil {
			return err
		}

		return createOrder(profile, tableName, username, shipAddress, itemsId)
	},
}

func init() {
	createOrderCmd.Flags().StringP("username", "u", "", "The username to save")
	createOrderCmd.MarkFlagRequired("username")

	createOrderCmd.Flags().StringP("ship-address", "a", "", "The address id to ship the order")
	createOrderCmd.MarkFlagRequired("ship-address")

	createOrderCmd.Flags().IntSliceP("items", "i", nil, "The items id as comma list to add [1-10]")
	createOrderCmd.MarkFlagRequired("items")

	rootCmd.AddCommand(createOrderCmd)
}

func createOrder(profile, tableName, username, shipAddress string, itemIds []int) error {
	repo, err := repository.NewDynamoDBRepo(profile, tableName)
	if err != nil {
		return err
	}

	customer, err := repo.GetCustomer(context.Background(), username)
	if err != nil {
		return err
	}
	customer.Addresses = map[string]repository.Address{shipAddress: {}}

	itemStock := repository.LoadItemStock()
	items := make(repository.OrderItems, len(itemIds))
	for i, itemId := range itemIds {
		item := itemStock.Get(itemId)
		items[i] = item
	}

	return repo.CreateOrder(context.Background(), customer, items)
}
