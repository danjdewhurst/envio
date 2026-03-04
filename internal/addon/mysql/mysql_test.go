package mysql

import "testing"

func TestMySQLInterface(t *testing.T) {
	m := New()
	if m.Name() != "mysql" {
		t.Errorf("expected name mysql, got %s", m.Name())
	}
	if m.DisplayName() != "MySQL" {
		t.Errorf("expected display name MySQL, got %s", m.DisplayName())
	}
}

func TestMySQLServices(t *testing.T) {
	m := New()
	services := m.Services()
	if len(services) != 1 || services[0].Name != "mysql" {
		t.Errorf("expected 1 mysql service, got %v", services)
	}
	if services[0].Image != "mysql:8.0" {
		t.Errorf("expected image mysql:8.0, got %s", services[0].Image)
	}
}

func TestMySQLVolumes(t *testing.T) {
	m := New()
	volumes := m.Volumes()
	if len(volumes) != 1 || volumes[0].Name != "mysql_data" {
		t.Errorf("expected mysql_data volume, got %v", volumes)
	}
}

func TestMySQLEnvVars(t *testing.T) {
	m := New()
	env := m.EnvVars()
	if env["DB_CONNECTION"] != "mysql" {
		t.Errorf("expected DB_CONNECTION=mysql, got %s", env["DB_CONNECTION"])
	}
	if env["DB_HOST"] != "mysql" {
		t.Errorf("expected DB_HOST=mysql, got %s", env["DB_HOST"])
	}
}
