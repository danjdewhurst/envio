package registry

import (
	"github.com/danjdewhurst/envio/internal/addon/mariadb"
	"github.com/danjdewhurst/envio/internal/addon/meilisearch"
	"github.com/danjdewhurst/envio/internal/addon/mysql"
	"github.com/danjdewhurst/envio/internal/addon/postgres"
	"github.com/danjdewhurst/envio/internal/addon/redis"
	"github.com/danjdewhurst/envio/internal/app/laravel"
)

// Default returns a registry pre-loaded with all built-in apps and addons.
func Default() *Registry {
	r := New()

	// Apps
	r.RegisterApp(laravel.New())

	// Addons
	r.RegisterAddon(mariadb.New())
	r.RegisterAddon(mysql.New())
	r.RegisterAddon(postgres.New())
	r.RegisterAddon(redis.New())
	r.RegisterAddon(meilisearch.New())

	return r
}
