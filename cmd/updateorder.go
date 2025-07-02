package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

func newUpdateStatusOrderCmd(app *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "update-status-order",
		Aliases:      []string{"uso"},
		Short:        "Update the order status",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			orderId, _ := cmd.Flags().GetString("order-id")
			status, _ := cmd.Flags().GetString("status")

			return updateStatusOrder(app.repo, orderId, status)
		},
	}

	cmd.Flags().StringP("order-id", "o", "", "The order id")
	cmd.MarkFlagRequired("order-id")

	cmd.Flags().StringP("status", "s", "", "The new status for the order")
	cmd.MarkFlagRequired("status")

	return cmd
}

func updateStatusOrder(repo Repository, orderId, status string) error {
	order, err := repo.GetOrderById(context.Background(), orderId)
	if err != nil {
		return err
	}

	return repo.UpdateStatusOrder(context.Background(), order, status)
}
