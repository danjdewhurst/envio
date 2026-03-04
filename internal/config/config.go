package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const configFile = "envio.yaml"

// ProjectConfig represents the saved project configuration.
type ProjectConfig struct {
	App    string   `yaml:"app"`
	Addons []string `yaml:"addons,omitempty"`
}

// Load reads the project config from the given directory.
func Load(dir string) (*ProjectConfig, error) {
	data, err := os.ReadFile(filepath.Join(dir, configFile))
	if err != nil {
		return nil, err
	}

	var cfg ProjectConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// Save writes the project config to the given directory.
func Save(dir string, cfg *ProjectConfig) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, configFile), data, 0644)
}

// Exists checks if a project config exists in the given directory.
func Exists(dir string) bool {
	_, err := os.Stat(filepath.Join(dir, configFile))
	return err == nil
}
