package cli

import (
	"github.com/spf13/cobra"

	"github.com/gustavooferreira/prometheus-metrics-generator/internal/config"
)

func NewRootCmd() *cobra.Command {
	// rootCmd represents the base command when called without any subcommands
	rootCmd := &cobra.Command{
		Use:   "promgen",
		Short: "A cli tool for generating metric dynamically",
		Long: `A cli tool for generating metric dynamically.

This tool currently supports getting access tokens using the passwordless oauth flow from Auth0.`,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	defaultConfigPath := config.GetDefaultConfigPath()

	// Global persistent flags
	rootCmd.PersistentFlags().StringP("config", "c", defaultConfigPath, "config file path")

	// Init and register sub commands
	_ = newVersionCmd(rootCmd)

	return rootCmd
}
