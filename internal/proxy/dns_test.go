package proxy

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteComposeFileCreatesDirectory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nested", "dir")

	if err := WriteComposeFile(dir); err != nil {
		t.Fatalf("WriteComposeFile failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(dir, "docker-compose.yml")); err != nil {
		t.Error("docker-compose.yml not created in nested directory")
	}
}

func TestHostsMarkers(t *testing.T) {
	if hostsMarkerBegin == "" || hostsMarkerEnd == "" {
		t.Error("hosts markers should not be empty")
	}
}
