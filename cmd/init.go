package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/danjdewhurst/envio/internal/app"
	"github.com/danjdewhurst/envio/internal/compose"
	"github.com/danjdewhurst/envio/internal/config"
)

var initAddons []string
var initVariant string

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

		selectedApp, err := reg.GetApp(appName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unknown app: %s\n\nAvailable apps:\n", appName)
			for _, a := range reg.ListApps() {
				fmt.Fprintf(os.Stderr, "  - %s: %s\n", a.Name(), a.Description())
			}
			return fmt.Errorf("unknown app: %s", appName)
		}

		// Apply variant if specified
		if initVariant != "" {
			va, ok := selectedApp.(app.VariantApp)
			if !ok {
				return fmt.Errorf("app %q does not support variants", appName)
			}
			if err := va.SetVariant(initVariant); err != nil {
				return err
			}
		}

		// Validate addons
		available := selectedApp.AvailableAddons()
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
				return fmt.Errorf("addon %q is not compatible with %s", addonName, selectedApp.DisplayName())
			}
		}

		// Build compose file
		cf := compose.NewComposeFile()
		cf.AddNetwork(compose.Network{Name: "envio", Driver: "bridge"})

		for _, svc := range selectedApp.Services() {
			cf.AddService(svc)
		}
		for _, vol := range selectedApp.Volumes() {
			cf.AddVolume(vol)
		}

		// Collect environment variables from app defaults and addons
		env := make(map[string]string)
		for k, v := range selectedApp.DefaultEnv() {
			env[k] = v
		}

		for _, addonName := range initAddons {
			addon, _ := reg.GetAddon(addonName)
			for _, svc := range addon.Services() {
				cf.AddService(svc)
			}
			for _, vol := range addon.Volumes() {
				cf.AddVolume(vol)
			}
			for k, v := range addon.EnvVars() {
				env[k] = v
			}
		}

		// Inject collected env vars into the app service
		if len(env) > 0 {
			if svc, ok := cf.Services["app"]; ok {
				if svc.Environment == nil {
					svc.Environment = make(map[string]string)
				}
				for k, v := range env {
					svc.Environment[k] = v
				}
				cf.Services["app"] = svc
			}
		}

		// Generate docker-compose.yml
		if err := compose.Generate(dir, cf); err != nil {
			return fmt.Errorf("failed to generate docker-compose.yml: %w", err)
		}

		// Write scaffold files if the app provides them
		var scaffoldPaths []string
		if scaffolder, ok := selectedApp.(app.Scaffolder); ok {
			for _, f := range scaffolder.ScaffoldFiles() {
				fullPath := filepath.Join(dir, f.Path)
				if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
					return fmt.Errorf("failed to create directory for %s: %w", f.Path, err)
				}
				if err := os.WriteFile(fullPath, []byte(f.Content), 0644); err != nil {
					return fmt.Errorf("failed to write %s: %w", f.Path, err)
				}
				scaffoldPaths = append(scaffoldPaths, f.Path)
			}
		}

		// Save envio config
		cfg := &config.ProjectConfig{
			App:     appName,
			Variant: initVariant,
			Addons:  initAddons,
		}
		if err := config.Save(dir, cfg); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Printf("Initialized %s environment", selectedApp.DisplayName())
		if initVariant != "" {
			fmt.Printf(" (variant: %s)", initVariant)
		}
		if len(initAddons) > 0 {
			fmt.Printf(" with addons: %v", initAddons)
		}
		fmt.Println()
		fmt.Println("\nGenerated files:")
		fmt.Println("  - docker-compose.yml")
		fmt.Println("  - envio.yaml")
		for _, p := range scaffoldPaths {
			fmt.Printf("  - %s\n", p)
		}
		fmt.Println("\nRun 'envio up' to start your environment.")

		return nil
	},
}

func init() {
	initCmd.Flags().StringSliceVarP(&initAddons, "addon", "a", nil, "Addons to include (e.g. --addon redis --addon mysql)")
	initCmd.Flags().StringVarP(&initVariant, "variant", "v", "", "App variant to use (e.g. --variant frankenphp)")
	rootCmd.AddCommand(initCmd)
}
