package postgres

import "github.com/danjdewhurst/envio/internal/compose"

type Postgres struct{}

func New() *Postgres { return &Postgres{} }

func (p *Postgres) Name() string        { return "postgres" }
func (p *Postgres) DisplayName() string  { return "PostgreSQL" }
func (p *Postgres) Description() string  { return "PostgreSQL 16 database server" }

func (p *Postgres) Services() []compose.Service {
	return []compose.Service{
		{
			Name:  "postgres",
			Image: "postgres:16-alpine",
			Ports: []string{"5432:5432"},
			Environment: map[string]string{
				"POSTGRES_DB":       "${DB_DATABASE:-envio}",
				"POSTGRES_USER":     "${DB_USERNAME:-envio}",
				"POSTGRES_PASSWORD": "${DB_PASSWORD:-secret}",
			},
			Volumes:  []string{"postgres_data:/var/lib/postgresql/data"},
			Networks: []string{"envio"},
			Restart:  "unless-stopped",
		},
	}
}

func (p *Postgres) Volumes() []compose.Volume {
	return []compose.Volume{
		{Name: "postgres_data"},
	}
}

func (p *Postgres) EnvVars() map[string]string {
	return map[string]string{
		"DB_CONNECTION": "pgsql",
		"DB_HOST":       "postgres",
		"DB_PORT":       "5432",
		"DB_DATABASE":   "envio",
		"DB_USERNAME":   "envio",
		"DB_PASSWORD":   "secret",
	}
}
