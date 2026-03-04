package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/danjdewhurst/envio/internal/config"
	"github.com/danjdewhurst/envio/internal/proxy"
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

		// Warn if project has a domain but proxy isn't running
		cfg, err := config.Load(dir)
		if err == nil && cfg.Domain != "" && !proxy.IsRunning() {
			fmt.Println("Warning: project has domain configured but proxy is not running.")
			fmt.Printf("Run 'envio proxy start' to access http://%s.test\n", cfg.Domain)
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
