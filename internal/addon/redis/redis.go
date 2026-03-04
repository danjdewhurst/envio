package redis

import "github.com/danjdewhurst/envio/internal/compose"

type Redis struct{}

func New() *Redis { return &Redis{} }

func (r *Redis) Name() string        { return "redis" }
func (r *Redis) DisplayName() string  { return "Redis" }
func (r *Redis) Description() string  { return "Redis in-memory data store" }

func (r *Redis) Services() []compose.Service {
	return []compose.Service{
		{
			Name:     "redis",
			Image:    "redis:alpine",
			Ports:    []string{"6379:6379"},
			Volumes:  []string{"redis_data:/data"},
			Networks: []string{"envio"},
			Restart:  "unless-stopped",
		},
	}
}

func (r *Redis) Volumes() []compose.Volume {
	return []compose.Volume{
		{Name: "redis_data"},
	}
}

func (r *Redis) EnvVars() map[string]string {
	return map[string]string{
		"REDIS_HOST": "redis",
		"REDIS_PORT": "6379",
	}
}
