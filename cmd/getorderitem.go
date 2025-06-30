package cmd

import (
	"context"
	"fmt"
	"meetup/repository"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var getOrderItemsCmd = &cobra.Command{
	Use:          "get-order-items",
	Aliases:      []string{"gi"},
	Short:        "Get the order's items",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		profile := viper.GetString("profile")
		tableName := viper.GetString("table")

		orderId, err := cmd.Flags().GetString("order-id")
		if err != nil {
			return err
		}

		return getOrderItems(profile, tableName, orderId)
	},
}

func init() {
	getOrderItemsCmd.Flags().StringP("order-id", "o", "", "The order id")
	getOrderItemsCmd.MarkFlagRequired("order-id")

	// a pesquisa será feita no GSI1, a PK será o id da order e o sk começa com ITEM#

	rootCmd.AddCommand(getOrderItemsCmd)
}

func getOrderItems(profile, tableName, orderId string) error {
	repo, err := repository.NewDynamoDBRepo(profile, tableName)
	if err != nil {
		return err
	}

	orderItems, err := repo.GetOrderItems(context.Background(), orderId)
	if err != nil {
		return err
	}

	return printOrderItems(orderItems)
}

func printOrderItems(orderItems repository.OrderItems) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

	fmt.Fprintln(w, "#\tId\tDescription\tPrice\t")
	for k, v := range orderItems {
		price := float32(v.Price) / 100

		fmt.Fprintf(w, "%d\t%s\t%s\t%.2f\t\n", k+1, v.Id, v.Description, price)
	}

	return w.Flush()
}
