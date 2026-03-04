package cmd

import (
	"fmt"
	"os"

	"github.com/danjdewhurst/envio/internal/compose"
	"github.com/danjdewhurst/envio/internal/config"
	"github.com/spf13/cobra"
)

var initAddons []string

var initCmd = &cobra.Command{
	Use:   "init [app]",
	Short: "Initialize a new Docker environment",
	Long:  `Initialize a new Docker development environment for the specified application type.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		appName := args[0]

		dir, err := os.Getwd()
		if err != nil {
			return err
		}

		if config.Exists(dir) {
			return fmt.Errorf("envio project already exists in this directory")
		}

		app, err := reg.GetApp(appName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unknown app: %s\n\nAvailable apps:\n", appName)
			for _, a := range reg.ListApps() {
				fmt.Fprintf(os.Stderr, "  - %s: %s\n", a.Name(), a.Description())
			}
			return fmt.Errorf("unknown app: %s", appName)
		}

		// Validate addons
		available := app.AvailableAddons()
		for _, addonName := range initAddons {
			if _, err := reg.GetAddon(addonName); err != nil {
				return err
			}
			found := false
			for _, a := range available {
				if a == addonName {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("addon %q is not compatible with %s", addonName, app.DisplayName())
			}
		}

		// Build compose file
		cf := compose.NewComposeFile()
		cf.AddNetwork(compose.Network{Name: "envio", Driver: "bridge"})

		for _, svc := range app.Services() {
			cf.AddService(svc)
		}
		for _, vol := range app.Volumes() {
			cf.AddVolume(vol)
		}

		for _, addonName := range initAddons {
			addon, _ := reg.GetAddon(addonName)
			for _, svc := range addon.Services() {
				cf.AddService(svc)
			}
			for _, vol := range addon.Volumes() {
				cf.AddVolume(vol)
			}
		}

		// Generate docker-compose.yml
		if err := compose.Generate(dir, cf); err != nil {
			return fmt.Errorf("failed to generate docker-compose.yml: %w", err)
		}

		// Save envio config
		cfg := &config.ProjectConfig{
			App:    appName,
			Addons: initAddons,
		}
		if err := config.Save(dir, cfg); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Printf("Initialized %s environment", app.DisplayName())
		if len(initAddons) > 0 {
			fmt.Printf(" with addons: %v", initAddons)
		}
		fmt.Println()
		fmt.Println("\nGenerated files:")
		fmt.Println("  - docker-compose.yml")
		fmt.Println("  - envio.yaml")
		fmt.Println("\nRun 'envio up' to start your environment.")

		return nil
	},
}

func init() {
	initCmd.Flags().StringSliceVarP(&initAddons, "addon", "a", nil, "Addons to include (e.g. --addon redis --addon mysql)")
	rootCmd.AddCommand(initCmd)
}
