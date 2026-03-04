package cmd

import (
	"github.com/danjdewhurst/envio/internal/registry"
	"github.com/spf13/cobra"
)

var reg *registry.Registry

var rootCmd = &cobra.Command{
	Use:   "envio",
	Short: "Create and manage local Docker setups for your applications",
	Long:  `Envio is a CLI tool that helps you create and manage local Docker development environments for various application types like Laravel, WordPress, Next.js and more.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	reg = registry.Default()
}
