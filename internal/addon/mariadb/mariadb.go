package mariadb

import "github.com/danjdewhurst/envio/internal/compose"

type MariaDB struct{}

func New() *MariaDB { return &MariaDB{} }

func (m *MariaDB) Name() string        { return "mariadb" }
func (m *MariaDB) DisplayName() string { return "MariaDB" }
func (m *MariaDB) Description() string { return "MariaDB 11.4 LTS database server" }

func (m *MariaDB) Services() []compose.Service {
	return []compose.Service{
		{
			Name:  "mariadb",
			Image: "mariadb:11.4",
			Ports: []string{"3306:3306"},
			Environment: map[string]string{
				"MARIADB_ROOT_PASSWORD": "secret",
				"MARIADB_DATABASE":      "envio",
				"MARIADB_USER":          "envio",
				"MARIADB_PASSWORD":      "secret",
			},
			Volumes:  []string{"mariadb_data:/var/lib/mysql"},
			Networks: []string{"envio"},
			Restart:  "unless-stopped",
		},
	}
}

func (m *MariaDB) Volumes() []compose.Volume {
	return []compose.Volume{
		{Name: "mariadb_data"},
	}
}

func (m *MariaDB) EnvVars() map[string]string {
	return map[string]string{
		"DB_CONNECTION": "mysql",
		"DB_HOST":       "mariadb",
		"DB_PORT":       "3306",
		"DB_DATABASE":   "envio",
		"DB_USERNAME":   "envio",
		"DB_PASSWORD":   "secret",
	}
}
