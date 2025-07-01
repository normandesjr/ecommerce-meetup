package cmd

import (
	"context"
	"meetup/repository"

	"github.com/spf13/cobra"
)

func newCreateOrderCmd(app *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "create-order",
		Aliases:      []string{"co"},
		Short:        "Create an order for the choosed items",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			username, _ := cmd.Flags().GetString("username")
			shipAddress, _ := cmd.Flags().GetString("ship-address")
			itemsId, _ := cmd.Flags().GetIntSlice("items")

			return createOrder(app.repo, username, shipAddress, itemsId)
		},
	}
	cmd.Flags().StringP("username", "u", "", "The username used to create the order")
	cmd.MarkFlagRequired("username")

	cmd.Flags().StringP("ship-address", "a", "", "The address id to ship the order")
	cmd.MarkFlagRequired("ship-address")

	cmd.Flags().IntSliceP("items", "i", nil, "The items id as comma list to add [1-10]")
	cmd.MarkFlagRequired("items")

	return cmd
}

func createOrder(repo Repository, username, shipAddress string, itemIds []int) error {
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
