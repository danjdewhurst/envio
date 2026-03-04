package redis

import "testing"

func TestRedisInterface(t *testing.T) {
	r := New()
	if r.Name() != "redis" {
		t.Errorf("expected name redis, got %s", r.Name())
	}
	if r.DisplayName() != "Redis" {
		t.Errorf("expected display name Redis, got %s", r.DisplayName())
	}
}

func TestRedisServices(t *testing.T) {
	r := New()
	services := r.Services()
	if len(services) != 1 || services[0].Name != "redis" {
		t.Errorf("expected 1 redis service, got %v", services)
	}
}

func TestRedisVolumes(t *testing.T) {
	r := New()
	volumes := r.Volumes()
	if len(volumes) != 1 || volumes[0].Name != "redis_data" {
		t.Errorf("expected redis_data volume, got %v", volumes)
	}
}

func TestRedisEnvVars(t *testing.T) {
	r := New()
	env := r.EnvVars()
	expected := map[string]string{
		"REDIS_HOST": "redis",
		"REDIS_PORT": "6379",
	}
	for key, want := range expected {
		if got := env[key]; got != want {
			t.Errorf("expected %s=%s, got %s", key, want, got)
		}
	}
}
