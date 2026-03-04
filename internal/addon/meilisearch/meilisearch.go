package meilisearch

import "github.com/danjdewhurst/envio/internal/compose"

type Meilisearch struct{}

func New() *Meilisearch { return &Meilisearch{} }

func (m *Meilisearch) Name() string        { return "meilisearch" }
func (m *Meilisearch) DisplayName() string  { return "Meilisearch" }
func (m *Meilisearch) Description() string  { return "Meilisearch search engine" }

func (m *Meilisearch) Services() []compose.Service {
	return []compose.Service{
		{
			Name:  "meilisearch",
			Image: "getmeili/meilisearch:latest",
			Ports: []string{"7700:7700"},
			Environment: map[string]string{
				"MEILI_NO_ANALYTICS": "true",
			},
			Volumes:  []string{"meilisearch_data:/meili_data"},
			Networks: []string{"envio"},
			Restart:  "unless-stopped",
		},
	}
}

func (m *Meilisearch) Volumes() []compose.Volume {
	return []compose.Volume{
		{Name: "meilisearch_data"},
	}
}

func (m *Meilisearch) EnvVars() map[string]string {
	return map[string]string{
		"SCOUT_DRIVER":       "meilisearch",
		"MEILISEARCH_HOST":   "http://meilisearch:7700",
		"MEILISEARCH_NO_ANALYTICS": "true",
	}
}
