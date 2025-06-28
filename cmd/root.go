package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "cloudday",
	Short: "CloudDay Triângulo",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringP("profile", "p", "", "AWS profile")

	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix("CDAY")
	viper.AutomaticEnv()

	viper.BindPFlag("profile", rootCmd.PersistentFlags().Lookup("profile"))
}

func initConfig() {
	profile := viper.GetString("profile")
	if profile == "" {
		fmt.Fprintln(os.Stderr, "Error: required flag \"profile\" not set")
		os.Exit(1)
	}
}
