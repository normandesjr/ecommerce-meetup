package cmd

import (
	"context"
	"fmt"
	"meetup/repository"
	"os"
	"text/tabwriter"
	"time"

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

	return print(orders)
}

func print(orders []repository.Order) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

	fmt.Fprintln(w, "#\tCreated at\tStatus\tShipped to\tTotal\tId\t")
	for k, v := range orders {
		createdAt, err := time.Parse(time.RFC3339, v.CreatedAt)
		if err != nil {
			return err
		}
		formattedDate := createdAt.Local().Format("02/01/2006")
		total := float32(v.Total) / 100

		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%.2f\t%s\t\n", k+1, formattedDate, v.Status, v.ShippedTo, total, v.Id)
	}

	return w.Flush()
}
