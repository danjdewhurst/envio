package postgres

import "testing"

func TestPostgresInterface(t *testing.T) {
	p := New()
	if p.Name() != "postgres" {
		t.Errorf("expected name postgres, got %s", p.Name())
	}
	if p.DisplayName() != "PostgreSQL" {
		t.Errorf("expected display name PostgreSQL, got %s", p.DisplayName())
	}
}

func TestPostgresServices(t *testing.T) {
	p := New()
	services := p.Services()
	if len(services) != 1 || services[0].Name != "postgres" {
		t.Errorf("expected 1 postgres service, got %v", services)
	}
	if services[0].Image != "postgres:16-alpine" {
		t.Errorf("expected image postgres:16-alpine, got %s", services[0].Image)
	}
}

func TestPostgresVolumes(t *testing.T) {
	p := New()
	volumes := p.Volumes()
	if len(volumes) != 1 || volumes[0].Name != "postgres_data" {
		t.Errorf("expected postgres_data volume, got %v", volumes)
	}
}

func TestPostgresEnvVars(t *testing.T) {
	p := New()
	env := p.EnvVars()
	expected := map[string]string{
		"DB_CONNECTION": "pgsql",
		"DB_HOST":       "postgres",
		"DB_PORT":       "5432",
		"DB_DATABASE":   "envio",
		"DB_USERNAME":   "envio",
		"DB_PASSWORD":   "secret",
	}
	for key, want := range expected {
		if got := env[key]; got != want {
			t.Errorf("expected %s=%s, got %s", key, want, got)
		}
	}
}
