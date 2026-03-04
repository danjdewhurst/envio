package registry

import (
	"fmt"

	"github.com/danjdewhurst/envio/internal/addon"
	"github.com/danjdewhurst/envio/internal/app"
)

// Registry holds all registered apps and addons.
type Registry struct {
	apps   map[string]app.App
	addons map[string]addon.Addon
}

// New creates an empty registry.
func New() *Registry {
	return &Registry{
		apps:   make(map[string]app.App),
		addons: make(map[string]addon.Addon),
	}
}

// RegisterApp adds an app to the registry.
func (r *Registry) RegisterApp(a app.App) {
	r.apps[a.Name()] = a
}

// RegisterAddon adds an addon to the registry.
func (r *Registry) RegisterAddon(a addon.Addon) {
	r.addons[a.Name()] = a
}

// GetApp returns a registered app by name.
func (r *Registry) GetApp(name string) (app.App, error) {
	a, ok := r.apps[name]
	if !ok {
		return nil, fmt.Errorf("unknown app: %s", name)
	}
	return a, nil
}

// GetAddon returns a registered addon by name.
func (r *Registry) GetAddon(name string) (addon.Addon, error) {
	a, ok := r.addons[name]
	if !ok {
		return nil, fmt.Errorf("unknown addon: %s", name)
	}
	return a, nil
}

// ListApps returns all registered app names.
func (r *Registry) ListApps() []app.App {
	apps := make([]app.App, 0, len(r.apps))
	for _, a := range r.apps {
		apps = append(apps, a)
	}
	return apps
}

// ListAddons returns all registered addon names.
func (r *Registry) ListAddons() []addon.Addon {
	addons := make([]addon.Addon, 0, len(r.addons))
	for _, a := range r.addons {
		addons = append(addons, a)
	}
	return addons
}
