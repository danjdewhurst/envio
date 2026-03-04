package compose

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Generate creates a docker-compose.yml file from the given compose file struct.
func Generate(dir string, cf *ComposeFile) error {
	data, err := yaml.Marshal(cf)
	if err != nil {
		return err
	}

	path := filepath.Join(dir, "docker-compose.yml")
	return os.WriteFile(path, data, 0644)
}

// NewComposeFile creates an empty ComposeFile with initialized maps.
func NewComposeFile() *ComposeFile {
	return &ComposeFile{
		Services: make(map[string]Service),
		Volumes:  make(map[string]Volume),
		Networks: make(map[string]Network),
	}
}

// AddService adds a service to the compose file.
func (cf *ComposeFile) AddService(s Service) {
	cf.Services[s.Name] = s
}

// AddVolume adds a named volume to the compose file.
func (cf *ComposeFile) AddVolume(v Volume) {
	cf.Volumes[v.Name] = v
}

// AddNetwork adds a network to the compose file.
func (cf *ComposeFile) AddNetwork(n Network) {
	cf.Networks[n.Name] = n
}
