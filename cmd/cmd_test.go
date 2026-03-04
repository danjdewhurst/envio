package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitCommand(t *testing.T) {
	dir := t.TempDir()
	originalDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(originalDir)

	rootCmd.SetArgs([]string{"init", "laravel", "--addon", "redis"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("init command failed: %v", err)
	}

	// Check docker-compose.yml was created
	if _, err := os.Stat(filepath.Join(dir, "docker-compose.yml")); err != nil {
		t.Error("docker-compose.yml not created")
	}

	// Check envio.yaml was created
	if _, err := os.Stat(filepath.Join(dir, "envio.yaml")); err != nil {
		t.Error("envio.yaml not created")
	}
}

func TestInitCommandUnknownApp(t *testing.T) {
	dir := t.TempDir()
	originalDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(originalDir)

	rootCmd.SetArgs([]string{"init", "nonexistent"})
	if err := rootCmd.Execute(); err == nil {
		t.Error("init should fail for unknown app")
	}
}

func TestInitCommandInvalidAddon(t *testing.T) {
	dir := t.TempDir()
	originalDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(originalDir)

	rootCmd.SetArgs([]string{"init", "laravel", "--addon", "nonexistent"})
	if err := rootCmd.Execute(); err == nil {
		t.Error("init should fail for unknown addon")
	}
}

func TestInitCommandDuplicateProject(t *testing.T) {
	dir := t.TempDir()
	originalDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(originalDir)

	// Reset flags from previous tests
	initAddons = nil

	// First init should succeed
	rootCmd.SetArgs([]string{"init", "laravel"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("first init failed: %v", err)
	}

	// Second init should fail
	rootCmd.SetArgs([]string{"init", "laravel"})
	if err := rootCmd.Execute(); err == nil {
		t.Error("second init should fail with existing project")
	}
}

func TestAppsCommand(t *testing.T) {
	rootCmd.SetArgs([]string{"apps"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("apps command failed: %v", err)
	}
}

func TestAddonsCommand(t *testing.T) {
	rootCmd.SetArgs([]string{"addons"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("addons command failed: %v", err)
	}
}
