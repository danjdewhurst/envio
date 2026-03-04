package laravel

import "github.com/danjdewhurst/envio/internal/compose"

// Laravel implements the App interface for Laravel projects.
type Laravel struct{}

func New() *Laravel { return &Laravel{} }

func (l *Laravel) Name() string        { return "laravel" }
func (l *Laravel) DisplayName() string  { return "Laravel" }
func (l *Laravel) Description() string  { return "PHP Laravel application with Nginx and PHP-FPM" }

func (l *Laravel) Services() []compose.Service {
	return []compose.Service{
		{
			Name:  "app",
			Image: "php:8.3-fpm",
			Volumes: []string{
				".:/var/www/html",
			},
			Networks: []string{"envio"},
			Restart:  "unless-stopped",
		},
		{
			Name:  "web",
			Image: "nginx:alpine",
			Ports: []string{"80:80"},
			Volumes: []string{
				".:/var/www/html",
				"./docker/nginx/default.conf:/etc/nginx/conf.d/default.conf",
			},
			DependsOn: []string{"app"},
			Networks:  []string{"envio"},
			Restart:   "unless-stopped",
		},
	}
}

func (l *Laravel) Volumes() []compose.Volume {
	return nil
}

func (l *Laravel) DefaultEnv() map[string]string {
	return map[string]string{
		"APP_ENV":   "local",
		"APP_DEBUG": "true",
		"APP_PORT":  "80",
	}
}

func (l *Laravel) AvailableAddons() []string {
	return []string{"mysql", "postgres", "redis", "meilisearch"}
}
