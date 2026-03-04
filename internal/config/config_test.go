package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()

	cfg := &ProjectConfig{
		App:    "laravel",
		Addons: []string{"mysql", "redis"},
	}

	if err := Save(dir, cfg); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := Load(dir)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.App != "laravel" {
		t.Errorf("expected app laravel, got %s", loaded.App)
	}
	if len(loaded.Addons) != 2 || loaded.Addons[0] != "mysql" || loaded.Addons[1] != "redis" {
		t.Errorf("unexpected addons: %v", loaded.Addons)
	}
}

func TestExists(t *testing.T) {
	dir := t.TempDir()

	if Exists(dir) {
		t.Error("Exists should return false for empty directory")
	}

	if err := os.WriteFile(filepath.Join(dir, "envio.yaml"), []byte("app: test\n"), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	if !Exists(dir) {
		t.Error("Exists should return true after creating config")
	}
}

func TestLoadMissing(t *testing.T) {
	dir := t.TempDir()
	_, err := Load(dir)
	if err == nil {
		t.Error("Load should fail for missing config")
	}
}

func TestSaveAndLoadVariant(t *testing.T) {
	dir := t.TempDir()

	cfg := &ProjectConfig{
		App:     "laravel",
		Variant: "frankenphp",
		Addons:  []string{"redis"},
	}

	if err := Save(dir, cfg); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := Load(dir)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.Variant != "frankenphp" {
		t.Errorf("expected variant frankenphp, got %s", loaded.Variant)
	}
	if loaded.App != "laravel" {
		t.Errorf("expected app laravel, got %s", loaded.App)
	}
}

func TestLoadMalformed(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "envio.yaml"), []byte("{{invalid yaml"), 0644); err != nil {
		t.Fatalf("failed to write malformed config: %v", err)
	}
	_, err := Load(dir)
	if err == nil {
		t.Error("Load should fail for malformed YAML")
	}
}

func TestSaveNoAddons(t *testing.T) {
	dir := t.TempDir()

	cfg := &ProjectConfig{App: "laravel"}
	if err := Save(dir, cfg); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := Load(dir)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if len(loaded.Addons) != 0 {
		t.Errorf("expected no addons, got %v", loaded.Addons)
	}
}
