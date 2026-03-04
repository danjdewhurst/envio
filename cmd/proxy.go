package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/danjdewhurst/envio/internal/proxy"
)

var proxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "Manage the shared Traefik reverse proxy for .test domains",
}

var proxyStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the Traefik reverse proxy",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Starting proxy...")
		if err := proxy.Start(); err != nil {
			return fmt.Errorf("failed to start proxy: %w", err)
		}
		fmt.Println("Proxy is running. Traefik is listening on port 80.")
		return nil
	},
}

var proxyStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the Traefik reverse proxy",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Stopping proxy...")
		if err := proxy.Stop(); err != nil {
			return fmt.Errorf("failed to stop proxy: %w", err)
		}
		fmt.Println("Proxy stopped.")
		return nil
	},
}

var proxyStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check if the proxy is running",
	Run: func(cmd *cobra.Command, args []string) {
		if proxy.IsRunning() {
			fmt.Println("Proxy is running.")
		} else {
			fmt.Println("Proxy is not running.")
		}
	},
}

var proxySetupDNSCmd = &cobra.Command{
	Use:   "setup-dns",
	Short: "Configure dnsmasq to resolve *.test domains (macOS only)",
	RunE: func(cmd *cobra.Command, args []string) error {
		return proxy.SetupDNS()
	},
}

func init() {
	proxyCmd.AddCommand(proxyStartCmd)
	proxyCmd.AddCommand(proxyStopCmd)
	proxyCmd.AddCommand(proxyStatusCmd)
	proxyCmd.AddCommand(proxySetupDNSCmd)
	rootCmd.AddCommand(proxyCmd)
}
