package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var addonsCmd = &cobra.Command{
	Use:   "addons",
	Short: "List available addons",
	Run: func(cmd *cobra.Command, args []string) {
		addons := reg.ListAddons()
		fmt.Println("Available addons:")
		for _, a := range addons {
			fmt.Printf("  %-15s %s\n", a.Name(), a.Description())
		}
	},
}

func init() {
	rootCmd.AddCommand(addonsCmd)
}
