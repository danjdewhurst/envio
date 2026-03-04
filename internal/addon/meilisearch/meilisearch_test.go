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

func TestMeilisearchVolumes(t *testing.T) {
	m := New()
	volumes := m.Volumes()
	if len(volumes) != 1 || volumes[0].Name != "meilisearch_data" {
		t.Errorf("expected meilisearch_data volume, got %v", volumes)
	}
}

func TestMeilisearchEnvVars(t *testing.T) {
	m := New()
	env := m.EnvVars()
	expected := map[string]string{
		"SCOUT_DRIVER":             "meilisearch",
		"MEILISEARCH_HOST":         "http://meilisearch:7700",
		"MEILISEARCH_NO_ANALYTICS": "true",
	}
	for key, want := range expected {
		if got := env[key]; got != want {
			t.Errorf("expected %s=%s, got %s", key, want, got)
		}
	}
}
