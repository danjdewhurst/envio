package laravel

import "testing"

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
