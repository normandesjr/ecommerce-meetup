package cmd

import (
	"context"
	"fmt"
	"meetup/repository"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

func newGetOrderItemsCmd(app *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "get-order-items",
		Aliases:      []string{"gi"},
		Short:        "Get the order's items",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			orderId, _ := cmd.Flags().GetString("order-id")

			return getOrderItems(app.repo, orderId)
		},
	}

	cmd.Flags().StringP("order-id", "o", "", "The order id")
	cmd.MarkFlagRequired("order-id")

	return cmd
}

func getOrderItems(repo Repository, orderId string) error {
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
