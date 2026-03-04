package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var appsCmd = &cobra.Command{
	Use:   "apps",
	Short: "List available application types",
	Run: func(cmd *cobra.Command, args []string) {
		apps := reg.ListApps()
		fmt.Println("Available apps:")
		for _, a := range apps {
			fmt.Printf("  %-15s %s\n", a.Name(), a.Description())
			addons := a.AvailableAddons()
			if len(addons) > 0 {
				fmt.Printf("  %15s Addons: %v\n", "", addons)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(appsCmd)
}
