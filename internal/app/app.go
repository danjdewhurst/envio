package app

import "github.com/danjdewhurst/envio/internal/compose"

// App defines the interface that all application types must implement.
// To add a new app (e.g. WordPress, Next.js), create a struct that
// satisfies this interface and register it in the registry.
type App interface {
	// Name returns the unique identifier for this app type.
	Name() string

	// DisplayName returns a human-friendly name for display.
	DisplayName() string

	// Description returns a short description of the app type.
	Description() string

	// Services returns the core Docker Compose services for this app.
	Services() []compose.Service

	// Volumes returns any named volumes required by the app.
	Volumes() []compose.Volume

	// DefaultEnv returns default environment variables for the app.
	DefaultEnv() map[string]string

	// AvailableAddons returns the names of addons compatible with this app.
	AvailableAddons() []string
}

// VariantApp is an optional interface that apps can implement to support
// variants (e.g. Laravel with FrankenPHP instead of Nginx+PHP-FPM).
type VariantApp interface {
	App
	// Variants returns the available variant names for this app.
	Variants() []string
	// SetVariant configures the app to use the given variant.
	SetVariant(variant string) error
}

// WebServicer is an optional interface that apps can implement to identify
// which service handles HTTP traffic (used for reverse proxy configuration).
type WebServicer interface {
	// WebServiceName returns the name of the service that serves HTTP traffic.
	WebServiceName() string
}

// ScaffoldFile represents a file to be created during project initialization.
type ScaffoldFile struct {
	Path    string
	Content string
}

// Scaffolder is an optional interface that apps can implement to provide
// additional files (e.g. Dockerfiles, config files) during initialization.
type Scaffolder interface {
	ScaffoldFiles() []ScaffoldFile
}
