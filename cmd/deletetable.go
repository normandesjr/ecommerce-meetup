package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

func newDeleteTable(app *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "delete-table",
		Aliases:      []string{"dt"},
		Short:        "Delete the DynamoDB table",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return deleteTable(app.repo)
		},
	}

	return cmd
}

func deleteTable(repo Repository) error {
	action := func() {
		fmt.Printf("...")
	}

	err := repo.DeleteTable(context.Background(), action)
	if err != nil {
		return err
	}

	fmt.Println("\nTable deleted")

	return nil
}
