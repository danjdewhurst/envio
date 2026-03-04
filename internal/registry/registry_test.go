package registry

import (
	"testing"

	"github.com/danjdewhurst/envio/internal/addon/mysql"
	"github.com/danjdewhurst/envio/internal/addon/redis"
	"github.com/danjdewhurst/envio/internal/app/laravel"
)

func TestRegisterAndGetApp(t *testing.T) {
	r := New()
	r.RegisterApp(laravel.New())

	app, err := r.GetApp("laravel")
	if err != nil {
		t.Fatalf("GetApp failed: %v", err)
	}
	if app.Name() != "laravel" {
		t.Errorf("expected laravel, got %s", app.Name())
	}
}

func TestGetAppUnknown(t *testing.T) {
	r := New()
	_, err := r.GetApp("nonexistent")
	if err == nil {
		t.Error("GetApp should fail for unknown app")
	}
}

func TestRegisterAndGetAddon(t *testing.T) {
	r := New()
	r.RegisterAddon(redis.New())

	addon, err := r.GetAddon("redis")
	if err != nil {
		t.Fatalf("GetAddon failed: %v", err)
	}
	if addon.Name() != "redis" {
		t.Errorf("expected redis, got %s", addon.Name())
	}
}

func TestGetAddonUnknown(t *testing.T) {
	r := New()
	_, err := r.GetAddon("nonexistent")
	if err == nil {
		t.Error("GetAddon should fail for unknown addon")
	}
}

func TestListApps(t *testing.T) {
	r := New()
	r.RegisterApp(laravel.New())

	apps := r.ListApps()
	if len(apps) != 1 {
		t.Fatalf("expected 1 app, got %d", len(apps))
	}
	if apps[0].Name() != "laravel" {
		t.Errorf("expected laravel, got %s", apps[0].Name())
	}
}

func TestListAddons(t *testing.T) {
	r := New()
	r.RegisterAddon(mysql.New())
	r.RegisterAddon(redis.New())

	addons := r.ListAddons()
	if len(addons) != 2 {
		t.Fatalf("expected 2 addons, got %d", len(addons))
	}
}

func TestDefault(t *testing.T) {
	r := Default()

	if _, err := r.GetApp("laravel"); err != nil {
		t.Error("Default registry missing laravel app")
	}
	for _, name := range []string{"mysql", "postgres", "redis", "meilisearch"} {
		if _, err := r.GetAddon(name); err != nil {
			t.Errorf("Default registry missing %s addon", name)
		}
	}
}
