package cmd

import (
	"context"
	"errors"
	"fmt"
	"meetup/repository"

	"github.com/spf13/cobra"
)

func newCreateTableCmd(app *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "create-table",
		Aliases:      []string{"ct"},
		Short:        "Create the DynamoDB table",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return createTable(app.repo)
		},
	}

	return cmd
}

func createTable(repo Repository) error {
	action := func() {
		fmt.Printf("...")
	}

	err := repo.CreateTable(context.Background(), action)
	if err != nil {
		if errors.Is(err, repository.ErrTableAlreadyExists) {
			fmt.Println("\nTable already exists")
			return nil
		}

		return err
	}

	fmt.Println("\nTable created")

	return nil
}
