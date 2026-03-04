package laravel

import (
	"strings"
	"testing"

	"github.com/danjdewhurst/envio/internal/app"
)

func TestLaravelInterface(t *testing.T) {
	l := New()

	if l.Name() != "laravel" {
		t.Errorf("expected name laravel, got %s", l.Name())
	}
	if l.DisplayName() != "Laravel" {
		t.Errorf("expected display name Laravel, got %s", l.DisplayName())
	}
	if l.Description() == "" {
		t.Error("description should not be empty")
	}
}

func TestLaravelServices(t *testing.T) {
	l := New()
	services := l.Services()

	if len(services) != 2 {
		t.Fatalf("expected 2 services, got %d", len(services))
	}

	names := map[string]bool{}
	for _, s := range services {
		names[s.Name] = true
	}
	if !names["app"] {
		t.Error("missing 'app' service")
	}
	if !names["web"] {
		t.Error("missing 'web' service")
	}
}

func TestLaravelDefaultServiceUsesBuild(t *testing.T) {
	l := New()
	services := l.Services()

	var appSvc *struct{ found bool }
	for _, s := range services {
		if s.Name == "app" {
			if s.Build == nil {
				t.Fatal("expected app service to use Build config, got nil")
			}
			if s.Build.Dockerfile != "docker/php/Dockerfile" {
				t.Errorf("expected Dockerfile docker/php/Dockerfile, got %s", s.Build.Dockerfile)
			}
			if s.Image != "" {
				t.Errorf("expected no Image when using Build, got %s", s.Image)
			}
			appSvc = &struct{ found bool }{true}
		}
	}
	if appSvc == nil {
		t.Fatal("app service not found")
	}
}

func TestLaravelAvailableAddons(t *testing.T) {
	l := New()
	addons := l.AvailableAddons()

	expected := map[string]bool{"mysql": true, "postgres": true, "redis": true, "meilisearch": true}
	for _, a := range addons {
		if !expected[a] {
			t.Errorf("unexpected addon: %s", a)
		}
		delete(expected, a)
	}
	for a := range expected {
		t.Errorf("missing addon: %s", a)
	}
}

func TestLaravelDefaultEnv(t *testing.T) {
	l := New()
	env := l.DefaultEnv()

	if env["APP_ENV"] != "local" {
		t.Errorf("expected APP_ENV=local, got %s", env["APP_ENV"])
	}
}

func TestLaravelVariants(t *testing.T) {
	l := New()
	variants := l.Variants()

	if len(variants) != 1 || variants[0] != "frankenphp" {
		t.Errorf("expected [frankenphp], got %v", variants)
	}
}

func TestLaravelSetVariantInvalid(t *testing.T) {
	l := New()
	if err := l.SetVariant("nonexistent"); err == nil {
		t.Error("expected error for invalid variant")
	}
}

func TestLaravelFrankenPHP(t *testing.T) {
	l := New()
	if err := l.SetVariant("frankenphp"); err != nil {
		t.Fatalf("SetVariant failed: %v", err)
	}

	// Should have a single service
	services := l.Services()
	if len(services) != 1 {
		t.Fatalf("expected 1 service for frankenphp, got %d", len(services))
	}

	svc := services[0]
	if svc.Name != "app" {
		t.Errorf("expected service name 'app', got %s", svc.Name)
	}
	if svc.Image != "dunglas/frankenphp:latest-php8.4" {
		t.Errorf("expected frankenphp php8.4 image, got %s", svc.Image)
	}

	// Should expose ports 80 and 443
	portMap := map[string]bool{}
	for _, p := range svc.Ports {
		portMap[p] = true
	}
	if !portMap["80:80"] {
		t.Error("missing port 80:80")
	}
	if !portMap["443:443"] {
		t.Error("missing port 443:443")
	}

	// Description should mention FrankenPHP
	if l.Description() != "PHP Laravel application with FrankenPHP" {
		t.Errorf("unexpected description: %s", l.Description())
	}

	// DefaultEnv should include SERVER_NAME
	env := l.DefaultEnv()
	if env["SERVER_NAME"] != ":80" {
		t.Errorf("expected SERVER_NAME=:80, got %s", env["SERVER_NAME"])
	}
}

func TestLaravelScaffoldFiles(t *testing.T) {
	l := New()

	// Verify it implements the Scaffolder interface
	var _ app.Scaffolder = l

	files := l.ScaffoldFiles()

	// Default variant should produce Dockerfile + opcache.ini
	if len(files) != 2 {
		t.Fatalf("expected 2 scaffold files, got %d", len(files))
	}

	pathMap := map[string]string{}
	for _, f := range files {
		pathMap[f.Path] = f.Content
	}

	// Check opcache.ini
	opcache, ok := pathMap["docker/php/opcache.ini"]
	if !ok {
		t.Fatal("missing docker/php/opcache.ini")
	}
	if !strings.Contains(opcache, "opcache.jit=1255") {
		t.Error("opcache.ini should enable JIT")
	}
	if !strings.Contains(opcache, "opcache.jit_buffer_size=128M") {
		t.Error("opcache.ini should set JIT buffer size")
	}

	// Check Dockerfile
	dockerfile, ok := pathMap["docker/php/Dockerfile"]
	if !ok {
		t.Fatal("missing docker/php/Dockerfile")
	}
	if !strings.Contains(dockerfile, "php:8.4-fpm") {
		t.Error("Dockerfile should use php:8.4-fpm")
	}
	if !strings.Contains(dockerfile, "opcache") {
		t.Error("Dockerfile should enable opcache")
	}
}

func TestLaravelFrankenPHPScaffoldFiles(t *testing.T) {
	l := New()
	if err := l.SetVariant("frankenphp"); err != nil {
		t.Fatalf("SetVariant failed: %v", err)
	}

	files := l.ScaffoldFiles()

	// FrankenPHP should only produce opcache.ini (no Dockerfile)
	if len(files) != 1 {
		t.Fatalf("expected 1 scaffold file for frankenphp, got %d", len(files))
	}
	if files[0].Path != "docker/php/opcache.ini" {
		t.Errorf("expected docker/php/opcache.ini, got %s", files[0].Path)
	}
}
