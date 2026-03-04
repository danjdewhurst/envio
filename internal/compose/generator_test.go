package compose

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewComposeFile(t *testing.T) {
	cf := NewComposeFile()
	if cf.Services == nil || cf.Volumes == nil || cf.Networks == nil {
		t.Fatal("NewComposeFile should initialize all maps")
	}
}

func TestAddService(t *testing.T) {
	cf := NewComposeFile()
	cf.AddService(Service{Name: "web", Image: "nginx:alpine"})

	svc, ok := cf.Services["web"]
	if !ok {
		t.Fatal("service 'web' not found")
	}
	if svc.Image != "nginx:alpine" {
		t.Errorf("expected image nginx:alpine, got %s", svc.Image)
	}
}

func TestAddVolume(t *testing.T) {
	cf := NewComposeFile()
	cf.AddVolume(Volume{Name: "data", Driver: "local"})

	vol, ok := cf.Volumes["data"]
	if !ok {
		t.Fatal("volume 'data' not found")
	}
	if vol.Driver != "local" {
		t.Errorf("expected driver local, got %s", vol.Driver)
	}
}

func TestAddNetwork(t *testing.T) {
	cf := NewComposeFile()
	cf.AddNetwork(Network{Name: "app", Driver: "bridge"})

	net, ok := cf.Networks["app"]
	if !ok {
		t.Fatal("network 'app' not found")
	}
	if net.Driver != "bridge" {
		t.Errorf("expected driver bridge, got %s", net.Driver)
	}
}

func TestGenerate(t *testing.T) {
	dir := t.TempDir()

	cf := NewComposeFile()
	cf.AddService(Service{
		Name:  "app",
		Image: "php:8.3-fpm",
		Ports: []string{"9000:9000"},
	})
	cf.AddVolume(Volume{Name: "app_data"})
	cf.AddNetwork(Network{Name: "envio", Driver: "bridge"})

	if err := Generate(dir, cf); err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	path := filepath.Join(dir, "docker-compose.yml")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read generated file: %v", err)
	}

	content := string(data)
	for _, want := range []string{"php:8.3-fpm", "9000:9000", "app_data", "bridge"} {
		if !contains(content, want) {
			t.Errorf("generated file missing %q", want)
		}
	}
}

func TestGenerateBadDirectory(t *testing.T) {
	cf := NewComposeFile()
	cf.AddService(Service{Name: "app", Image: "nginx:alpine"})

	err := Generate("/nonexistent/path/that/does/not/exist", cf)
	if err == nil {
		t.Error("Generate should fail for non-existent directory")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
