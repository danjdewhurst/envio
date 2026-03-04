package mysql

import "github.com/danjdewhurst/envio/internal/compose"

type MySQL struct{}

func New() *MySQL { return &MySQL{} }

func (m *MySQL) Name() string        { return "mysql" }
func (m *MySQL) DisplayName() string { return "MySQL" }
func (m *MySQL) Description() string { return "MySQL 8.0 database server" }

func (m *MySQL) Services() []compose.Service {
	return []compose.Service{
		{
			Name:  "mysql",
			Image: "mysql:8.0",
			Ports: []string{"3306:3306"},
			Environment: map[string]string{
				"MYSQL_ROOT_PASSWORD": "${DB_PASSWORD:-secret}",
				"MYSQL_DATABASE":      "${DB_DATABASE:-envio}",
				"MYSQL_USER":          "${DB_USERNAME:-envio}",
				"MYSQL_PASSWORD":      "${DB_PASSWORD:-secret}",
			},
			Volumes:  []string{"mysql_data:/var/lib/mysql"},
			Networks: []string{"envio"},
			Restart:  "unless-stopped",
		},
	}
}

func (m *MySQL) Volumes() []compose.Volume {
	return []compose.Volume{
		{Name: "mysql_data"},
	}
}

func (m *MySQL) EnvVars() map[string]string {
	return map[string]string{
		"DB_CONNECTION": "mysql",
		"DB_HOST":       "mysql",
		"DB_PORT":       "3306",
		"DB_DATABASE":   "envio",
		"DB_USERNAME":   "envio",
		"DB_PASSWORD":   "secret",
	}
}
