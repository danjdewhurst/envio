package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/danjdewhurst/envio/internal/config"
	"github.com/spf13/cobra"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Start the Docker environment",
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := os.Getwd()
		if err != nil {
			return err
		}

		if !config.Exists(dir) {
			return fmt.Errorf("no envio project found in this directory — run 'envio init' first")
		}

		fmt.Println("Starting environment...")
		c := exec.Command("docker", "compose", "up", "-d")
		c.Dir = dir
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		return c.Run()
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
}
