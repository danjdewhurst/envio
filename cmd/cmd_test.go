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

func TestInitCommandWithVariant(t *testing.T) {
	dir := t.TempDir()
	chdir(t, dir)

	// Reset flags from previous tests
	initAddons = nil
	initVariant = ""

	rootCmd.SetArgs([]string{"init", "laravel", "--variant", "frankenphp"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("init command with variant failed: %v", err)
	}

	// Check docker-compose.yml uses a Dockerfile build (not a direct image)
	data, err := os.ReadFile(filepath.Join(dir, "docker-compose.yml"))
	if err != nil {
		t.Fatalf("failed to read compose file: %v", err)
	}
	content := string(data)
	if !stringContains(content, "dockerfile: docker/php/Dockerfile") {
		t.Error("compose file should use Dockerfile build for frankenphp variant")
	}
	// Should NOT have a separate web/nginx service
	if stringContains(content, "nginx") {
		t.Error("compose file should not contain nginx when using frankenphp")
	}

	// Check envio.yaml has variant
	cfgData, err := os.ReadFile(filepath.Join(dir, "envio.yaml"))
	if err != nil {
		t.Fatalf("failed to read config: %v", err)
	}
	if !stringContains(string(cfgData), "variant: frankenphp") {
		t.Error("config should contain variant: frankenphp")
	}
}

func TestInitCommandInvalidVariant(t *testing.T) {
	dir := t.TempDir()
	chdir(t, dir)

	initAddons = nil
	initVariant = ""

	rootCmd.SetArgs([]string{"init", "laravel", "--variant", "nonexistent"})
	if err := rootCmd.Execute(); err == nil {
		t.Error("init should fail for invalid variant")
	}
}

func TestUpCommandNoConfig(t *testing.T) {
	dir := t.TempDir()
	chdir(t, dir)

	rootCmd.SetArgs([]string{"up"})
	if err := rootCmd.Execute(); err == nil {
		t.Error("up should fail without envio.yaml")
	}
}

func TestDownCommandNoConfig(t *testing.T) {
	dir := t.TempDir()
	chdir(t, dir)

	rootCmd.SetArgs([]string{"down"})
	if err := rootCmd.Execute(); err == nil {
		t.Error("down should fail without envio.yaml")
	}
}

func TestStatusCommandNoConfig(t *testing.T) {
	dir := t.TempDir()
	chdir(t, dir)

	rootCmd.SetArgs([]string{"status"})
	if err := rootCmd.Execute(); err == nil {
		t.Error("status should fail without envio.yaml")
	}
}

func TestInitCommandWithDomain(t *testing.T) {
	dir := t.TempDir()
	chdir(t, dir)

	initAddons = nil
	initVariant = ""
	initDomain = ""
	initNoProxy = false

	rootCmd.SetArgs([]string{"init", "laravel", "--domain", "myapp"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("init command with domain failed: %v", err)
	}

	// Check compose file has Traefik labels
	data, err := os.ReadFile(filepath.Join(dir, "docker-compose.yml"))
	if err != nil {
		t.Fatalf("failed to read compose file: %v", err)
	}
	content := string(data)

	if !stringContains(content, "traefik.enable") {
		t.Error("compose file missing traefik.enable label")
	}
	if !stringContains(content, "Host(`myapp.test`)") {
		t.Error("compose file missing Host rule")
	}
	if !stringContains(content, "envio-proxy") {
		t.Error("compose file missing envio-proxy network")
	}
	if !stringContains(content, "external: true") {
		t.Error("compose file missing external network")
	}

	// Web service should NOT have host port mappings
	if stringContains(content, "\"80:80\"") || stringContains(content, "- 80:80") {
		t.Error("web service should not have host port mappings when proxy is enabled")
	}

	// Check envio.yaml has domain
	cfgData, err := os.ReadFile(filepath.Join(dir, "envio.yaml"))
	if err != nil {
		t.Fatalf("failed to read config: %v", err)
	}
	if !stringContains(string(cfgData), "domain: myapp") {
		t.Error("config should contain domain: myapp")
	}
}

func TestInitCommandDefaultDomainFromDir(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "My Cool App")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("failed to create dir: %v", err)
	}
	chdir(t, dir)

	initAddons = nil
	initVariant = ""
	initDomain = ""
	initNoProxy = false

	rootCmd.SetArgs([]string{"init", "laravel"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("init command failed: %v", err)
	}

	// Default domain should be sanitized from dir name
	cfgData, err := os.ReadFile(filepath.Join(dir, "envio.yaml"))
	if err != nil {
		t.Fatalf("failed to read config: %v", err)
	}
	if !stringContains(string(cfgData), "domain: my-cool-app") {
		t.Errorf("config should contain domain: my-cool-app, got:\n%s", string(cfgData))
	}
}

func TestInitCommandNoProxy(t *testing.T) {
	dir := t.TempDir()
	chdir(t, dir)

	initAddons = nil
	initVariant = ""
	initDomain = ""
	initNoProxy = false

	rootCmd.SetArgs([]string{"init", "laravel", "--no-proxy"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("init command with --no-proxy failed: %v", err)
	}

	// Check compose file does NOT have Traefik labels
	data, err := os.ReadFile(filepath.Join(dir, "docker-compose.yml"))
	if err != nil {
		t.Fatalf("failed to read compose file: %v", err)
	}
	content := string(data)

	if stringContains(content, "traefik.enable") {
		t.Error("compose file should not have traefik labels with --no-proxy")
	}

	// Should still have port mappings
	if !stringContains(content, "80:80") {
		t.Error("compose file should preserve port mappings with --no-proxy")
	}

	// Config should have no domain
	cfgData, err := os.ReadFile(filepath.Join(dir, "envio.yaml"))
	if err != nil {
		t.Fatalf("failed to read config: %v", err)
	}
	if stringContains(string(cfgData), "domain:") {
		t.Error("config should not contain domain with --no-proxy")
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
