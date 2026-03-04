package mariadb

import "testing"

func TestMariaDBInterface(t *testing.T) {
	m := New()
	if m.Name() != "mariadb" {
		t.Errorf("expected name mariadb, got %s", m.Name())
	}
	if m.DisplayName() != "MariaDB" {
		t.Errorf("expected display name MariaDB, got %s", m.DisplayName())
	}
}

func TestMariaDBServices(t *testing.T) {
	m := New()
	services := m.Services()
	if len(services) != 1 || services[0].Name != "mariadb" {
		t.Errorf("expected 1 mariadb service, got %v", services)
	}
	if services[0].Image != "mariadb:11.4" {
		t.Errorf("expected image mariadb:11.4, got %s", services[0].Image)
	}
}

func TestMariaDBVolumes(t *testing.T) {
	m := New()
	volumes := m.Volumes()
	if len(volumes) != 1 || volumes[0].Name != "mariadb_data" {
		t.Errorf("expected mariadb_data volume, got %v", volumes)
	}
}

func TestMariaDBEnvVars(t *testing.T) {
	m := New()
	env := m.EnvVars()
	expected := map[string]string{
		"DB_CONNECTION": "mysql",
		"DB_HOST":       "mariadb",
		"DB_PORT":       "3306",
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
