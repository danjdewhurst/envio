package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/danjdewhurst/envio/internal/config"
	"github.com/spf13/cobra"
)

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Stop the Docker environment",
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := os.Getwd()
		if err != nil {
			return err
		}

		if !config.Exists(dir) {
			return fmt.Errorf("no envio project found in this directory")
		}

		fmt.Println("Stopping environment...")
		c := exec.Command("docker", "compose", "down")
		c.Dir = dir
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		return c.Run()
	},
}

func init() {
	rootCmd.AddCommand(downCmd)
}
