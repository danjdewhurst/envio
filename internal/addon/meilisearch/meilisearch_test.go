package meilisearch

import "testing"

func TestMeilisearchInterface(t *testing.T) {
	m := New()
	if m.Name() != "meilisearch" {
		t.Errorf("expected name meilisearch, got %s", m.Name())
	}
	if m.DisplayName() != "Meilisearch" {
		t.Errorf("expected display name Meilisearch, got %s", m.DisplayName())
	}
}

func TestMeilisearchServices(t *testing.T) {
	m := New()
	services := m.Services()
	if len(services) != 1 || services[0].Name != "meilisearch" {
		t.Errorf("expected 1 meilisearch service, got %v", services)
	}
}

func TestMeilisearchEnvVars(t *testing.T) {
	m := New()
	env := m.EnvVars()
	if env["SCOUT_DRIVER"] != "meilisearch" {
		t.Errorf("expected SCOUT_DRIVER=meilisearch, got %s", env["SCOUT_DRIVER"])
	}
}
