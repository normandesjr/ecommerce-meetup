package cmd

import (
	"context"
	"fmt"
	"meetup/repository"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var getOrderCmd = &cobra.Command{
	Use:          "get-order",
	Aliases:      []string{"go"},
	Short:        "Get the customer orders",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		profile := viper.GetString("profile")
		tableName := viper.GetString("table")

		username, err := cmd.Flags().GetString("username")
		if err != nil {
			return err
		}

		return getOrder(profile, tableName, username)
	},
}

func init() {
	getOrderCmd.Flags().StringP("username", "u", "", "The username to search the orders")
	getOrderCmd.MarkFlagRequired("username")

	rootCmd.AddCommand(getOrderCmd)
}

func getOrder(profile, tableName, username string) error {
	repo, err := repository.NewDynamoDBRepo(profile, tableName)
	if err != nil {
		return err
	}

	customer, err := repo.GetCustomer(context.Background(), username)
	if err != nil {
		return err
	}

	orders, err := repo.GetOrders(context.Background(), customer)
	if err != nil {
		return err
	}

	for _, order := range orders {
		fmt.Printf("%s: %s [%d] - %s\n", order.Id, order.Status, order.Total, order.ShippedTo)
	}

	return nil
}
