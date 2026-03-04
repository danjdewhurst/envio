package laravel

import (
	"fmt"

	"github.com/danjdewhurst/envio/internal/app"
	"github.com/danjdewhurst/envio/internal/compose"
)

// Laravel implements the App, VariantApp, and Scaffolder interfaces for Laravel projects.
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
				Image: "dunglas/frankenphp:latest-php8.4",
				Ports: []string{"80:80", "443:443"},
				Volumes: []string{
					".:/app",
					"./docker/php/opcache.ini:/usr/local/etc/php/conf.d/opcache.ini",
				},
				Networks: []string{"envio"},
				Restart:  "unless-stopped",
			},
		}
	}

	return []compose.Service{
		{
			Name: "app",
			Build: &compose.BuildConfig{
				Context:    ".",
				Dockerfile: "docker/php/Dockerfile",
			},
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

const phpDockerfile = `FROM php:8.4-fpm

# Install system dependencies
RUN apt-get update && apt-get install -y \
    git \
    curl \
    libpng-dev \
    libonig-dev \
    libxml2-dev \
    zip \
    unzip \
    && docker-php-ext-install pdo_mysql mbstring exif pcntl bcmath gd \
    && rm -rf /var/lib/apt/lists/*

# Enable and configure OPcache
RUN docker-php-ext-enable opcache

COPY docker/php/opcache.ini /usr/local/etc/php/conf.d/opcache.ini

# Install Composer
COPY --from=composer:latest /usr/bin/composer /usr/bin/composer

WORKDIR /var/www/html
`

const opcacheINI = `[opcache]
opcache.enable=1
opcache.memory_consumption=256
opcache.interned_strings_buffer=16
opcache.max_accelerated_files=20000
opcache.validate_timestamps=1
opcache.revalidate_freq=0
opcache.jit=1255
opcache.jit_buffer_size=128M
`

func (l *Laravel) ScaffoldFiles() []app.ScaffoldFile {
	files := []app.ScaffoldFile{
		{Path: "docker/php/opcache.ini", Content: opcacheINI},
	}
	if l.variant != "frankenphp" {
		files = append(files, app.ScaffoldFile{
			Path:    "docker/php/Dockerfile",
			Content: phpDockerfile,
		})
	}
	return files
}
