package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/danjdewhurst/envio/internal/config"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the status of the Docker environment",
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := os.Getwd()
		if err != nil {
			return err
		}

		if !config.Exists(dir) {
			return fmt.Errorf("no envio project found in this directory")
		}

		cfg, err := config.Load(dir)
		if err != nil {
			return err
		}

		fmt.Printf("App: %s\n", cfg.App)
		if len(cfg.Addons) > 0 {
			fmt.Printf("Addons: %v\n", cfg.Addons)
		}
		fmt.Println()

		c := exec.Command("docker", "compose", "ps")
		c.Dir = dir
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		return c.Run()
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
