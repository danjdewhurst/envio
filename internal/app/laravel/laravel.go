package laravel

import (
	"fmt"

	"github.com/danjdewhurst/envio/internal/compose"
)

// Laravel implements the App and VariantApp interfaces for Laravel projects.
type Laravel struct {
	variant string
}

func New() *Laravel { return &Laravel{} }

func (l *Laravel) Name() string        { return "laravel" }
func (l *Laravel) DisplayName() string { return "Laravel" }

func (l *Laravel) Description() string {
	if l.variant == "frankenphp" {
		return "PHP Laravel application with FrankenPHP"
	}
	return "PHP Laravel application with Nginx and PHP-FPM"
}

func (l *Laravel) Variants() []string {
	return []string{"frankenphp"}
}

func (l *Laravel) SetVariant(variant string) error {
	for _, v := range l.Variants() {
		if v == variant {
			l.variant = variant
			return nil
		}
	}
	return fmt.Errorf("unknown variant %q for %s (available: %v)", variant, l.DisplayName(), l.Variants())
}

func (l *Laravel) Services() []compose.Service {
	if l.variant == "frankenphp" {
		return []compose.Service{
			{
				Name:  "app",
				Image: "dunglas/frankenphp:latest-php8.3",
				Ports: []string{"80:80", "443:443"},
				Volumes: []string{
					".:/app",
				},
				Networks: []string{"envio"},
				Restart:  "unless-stopped",
			},
		}
	}

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
	env := map[string]string{
		"APP_ENV":   "local",
		"APP_DEBUG": "true",
		"APP_PORT":  "80",
	}
	if l.variant == "frankenphp" {
		env["SERVER_NAME"] = ":80"
	}
	return env
}

func (l *Laravel) AvailableAddons() []string {
	return []string{"mysql", "postgres", "redis", "meilisearch"}
}
