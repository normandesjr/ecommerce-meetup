package cmd

import (
	"context"
	"fmt"
	"meetup/repository"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Repository interface {
	CreateCustomer(ctx context.Context, customer repository.Customer) error
	GetCustomer(ctx context.Context, username string) (*repository.Customer, error)
	UpdateAddress(ctx context.Context, customer repository.Customer, address repository.Address) error
	RemoveAddress(ctx context.Context, customer repository.Customer, addressId string) error
	CreateOrder(ctx context.Context, customer *repository.Customer, items repository.OrderItems) error
	UpdateStatusOrder(ctx context.Context, order *repository.Order, status string) error
	GetOrders(ctx context.Context, customer *repository.Customer) ([]repository.Order, error)
	GetOrderById(ctx context.Context, orderId string) (*repository.Order, error)
	GetOrderItems(ctx context.Context, orderId string) (repository.OrderItems, error)
	CreateTable(ctx context.Context, action func()) error
	DeleteTable(ctx context.Context, action func()) error
}

type App struct {
	repo Repository
}

func Execute() {
	rootCmd := newRootCmd()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	app := &App{}
	cmd := &cobra.Command{
		Use:   "cloudday --profile <profile> --table <table name>",
		Short: "CloudDay Triângulo",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			profile := viper.GetString("profile")
			tableName := viper.GetString("table")

			awsCfg, err := config.LoadDefaultConfig(context.TODO(), func(lo *config.LoadOptions) error {
				lo.SharedConfigProfile = profile
				return nil
			})
			if err != nil {
				return fmt.Errorf("connecting to AWS using the profile %q: %v", profile, err)
			}

			client := dynamodb.NewFromConfig(awsCfg)
			repo, err := repository.NewDynamoDBRepo(client, tableName)
			if err != nil {
				return err
			}
			app.repo = repo
			return nil
		},
	}

	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix("CDAY")
	viper.AutomaticEnv()

	cmd.PersistentFlags().StringP("profile", "p", "", "AWS profile")
	cmd.PersistentFlags().StringP("table", "t", "CloudDayTable", "DynamoDB table name")

	viper.BindPFlag("profile", cmd.PersistentFlags().Lookup("profile"))
	viper.BindPFlag("table", cmd.PersistentFlags().Lookup("table"))

	cmd.AddCommand(newCreateCustomerCmd(app))
	cmd.AddCommand(newCreateOrderCmd(app))
	cmd.AddCommand(newCreateTableCmd(app))
	cmd.AddCommand(newDeleteTable(app))
	cmd.AddCommand(newGetOrderCmd(app))
	cmd.AddCommand(newGetOrderItemsCmd(app))
	cmd.AddCommand(newUpdateAddressCmd(app))
	cmd.AddCommand(newRemoveAddressCmd(app))
	cmd.AddCommand(newUpdateStatusOrderCmd(app))

	cobra.OnInitialize(initConfig)
	return cmd
}

func initConfig() {
	profile := viper.GetString("profile")
	if profile == "" {
		fmt.Fprintln(os.Stderr, "Error: required flag \"profile\" not set")
		os.Exit(1)
	}
}
