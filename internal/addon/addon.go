package addon

import "github.com/danjdewhurst/envio/internal/compose"

// Addon defines the interface for optional services that can be added
// to any compatible app (e.g. Redis, Meilisearch, PostgreSQL).
type Addon interface {
	// Name returns the unique identifier for this addon.
	Name() string

	// DisplayName returns a human-friendly name for display.
	DisplayName() string

	// Description returns a short description of the addon.
	Description() string

	// Services returns the Docker Compose services this addon provides.
	Services() []compose.Service

	// Volumes returns any named volumes required by the addon.
	Volumes() []compose.Volume

	// EnvVars returns environment variables to inject into the app service.
	EnvVars() map[string]string
}
