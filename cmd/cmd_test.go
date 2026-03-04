package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func chdir(t *testing.T, dir string) {
	t.Helper()
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("failed to chdir to %s: %v", dir, err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Fatalf("failed to restore working directory: %v", err)
		}
	})
}

func TestInitCommand(t *testing.T) {
	dir := t.TempDir()
	chdir(t, dir)

	rootCmd.SetArgs([]string{"init", "laravel", "--addon", "redis"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("init command failed: %v", err)
	}

	// Check docker-compose.yml was created
	composePath := filepath.Join(dir, "docker-compose.yml")
	if _, err := os.Stat(composePath); err != nil {
		t.Error("docker-compose.yml not created")
	}

	// Check envio.yaml was created
	if _, err := os.Stat(filepath.Join(dir, "envio.yaml")); err != nil {
		t.Error("envio.yaml not created")
	}

	// Check env vars are injected into compose file
	data, err := os.ReadFile(composePath)
	if err != nil {
		t.Fatalf("failed to read compose file: %v", err)
	}
	content := string(data)
	for _, want := range []string{"APP_ENV", "REDIS_HOST"} {
		if !stringContains(content, want) {
			t.Errorf("compose file missing env var %q", want)
		}
	}
}

func stringContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestInitCommandUnknownApp(t *testing.T) {
	dir := t.TempDir()
	chdir(t, dir)

	rootCmd.SetArgs([]string{"init", "nonexistent"})
	if err := rootCmd.Execute(); err == nil {
		t.Error("init should fail for unknown app")
	}
}

func TestInitCommandInvalidAddon(t *testing.T) {
	dir := t.TempDir()
	chdir(t, dir)

	rootCmd.SetArgs([]string{"init", "laravel", "--addon", "nonexistent"})
	if err := rootCmd.Execute(); err == nil {
		t.Error("init should fail for unknown addon")
	}
}

func TestInitCommandDuplicateProject(t *testing.T) {
	dir := t.TempDir()
	chdir(t, dir)

	// Reset flags from previous tests
	initAddons = nil

	// First init should succeed
	rootCmd.SetArgs([]string{"init", "laravel"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("first init failed: %v", err)
	}

	// Verify default env vars are present even without addons
	data, err := os.ReadFile(filepath.Join(dir, "docker-compose.yml"))
	if err != nil {
		t.Fatalf("failed to read compose file: %v", err)
	}
	if !stringContains(string(data), "APP_ENV") {
		t.Error("compose file missing default APP_ENV")
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
