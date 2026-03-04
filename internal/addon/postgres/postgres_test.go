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

func TestPostgresEnvVars(t *testing.T) {
	p := New()
	env := p.EnvVars()
	if env["DB_CONNECTION"] != "pgsql" {
		t.Errorf("expected DB_CONNECTION=pgsql, got %s", env["DB_CONNECTION"])
	}
}
