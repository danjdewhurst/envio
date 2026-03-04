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

func TestRedisEnvVars(t *testing.T) {
	r := New()
	env := r.EnvVars()
	if env["REDIS_HOST"] != "redis" {
		t.Errorf("expected REDIS_HOST=redis, got %s", env["REDIS_HOST"])
	}
}
