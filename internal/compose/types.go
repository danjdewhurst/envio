package compose

// Service represents a single service in a Docker Compose file.
type Service struct {
	Name        string            `yaml:"-"`
	Image       string            `yaml:"image,omitempty"`
	Build       *BuildConfig      `yaml:"build,omitempty"`
	Ports       []string          `yaml:"ports,omitempty"`
	Volumes     []string          `yaml:"volumes,omitempty"`
	Environment map[string]string `yaml:"environment,omitempty"`
	DependsOn   []string          `yaml:"depends_on,omitempty"`
	Networks    []string          `yaml:"networks,omitempty"`
	Command     string            `yaml:"command,omitempty"`
	Restart     string            `yaml:"restart,omitempty"`
	Labels      map[string]string `yaml:"labels,omitempty"`
}

// BuildConfig represents a Docker Compose build configuration.
type BuildConfig struct {
	Context    string `yaml:"context,omitempty"`
	Dockerfile string `yaml:"dockerfile,omitempty"`
}

// Volume represents a named volume in a Docker Compose file.
type Volume struct {
	Name   string `yaml:"-"`
	Driver string `yaml:"driver,omitempty"`
}

// Network represents a network in a Docker Compose file.
type Network struct {
	Name     string `yaml:"-"`
	Driver   string `yaml:"driver,omitempty"`
	External bool   `yaml:"external,omitempty"`
}

// ComposeFile represents a complete docker-compose.yml structure.
type ComposeFile struct {
	Services map[string]Service `yaml:"services"`
	Volumes  map[string]Volume  `yaml:"volumes,omitempty"`
	Networks map[string]Network `yaml:"networks,omitempty"`
}
