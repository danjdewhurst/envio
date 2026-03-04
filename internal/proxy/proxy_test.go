package proxy

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNetworkName(t *testing.T) {
	if name := NetworkName(); name != "envio-proxy" {
		t.Errorf("expected envio-proxy, got %s", name)
	}
}

func TestProxyDir(t *testing.T) {
	dir := ProxyDir()
	if !strings.HasSuffix(dir, filepath.Join(".envio", "proxy")) {
		t.Errorf("expected path ending in .envio/proxy, got %s", dir)
	}
}

func TestWriteComposeFile(t *testing.T) {
	dir := t.TempDir()

	if err := WriteComposeFile(dir); err != nil {
		t.Fatalf("WriteComposeFile failed: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(dir, "docker-compose.yml"))
	if err != nil {
		t.Fatalf("failed to read compose file: %v", err)
	}

	content := string(data)
	for _, want := range []string{
		"traefik:v3.3",
		"--providers.docker=true",
		"--providers.docker.exposedbydefault=false",
		"--providers.docker.network=envio-proxy",
		"--entrypoints.web.address=:80",
		"80:80",
		"/var/run/docker.sock:/var/run/docker.sock:ro",
		"envio-proxy",
		"external: true",
		"unless-stopped",
	} {
		if !strings.Contains(content, want) {
			t.Errorf("compose file missing %q", want)
		}
	}
}

func TestSanitiseDomain(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"My App", "my-app"},
		{"my_cool_app", "my-cool-app"},
		{"My--App", "my-app"},
		{"-leading-", "leading"},
		{"UPPERCASE", "uppercase"},
		{"hello world 123", "hello-world-123"},
		{"app@v2.0!", "app-v2-0"},
		{"---", ""},
	}

	for _, tt := range tests {
		got := SanitiseDomain(tt.input)
		if got != tt.want {
			t.Errorf("SanitiseDomain(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestTraefikLabels(t *testing.T) {
	labels := TraefikLabels("myapp", 80)

	expected := map[string]string{
		"traefik.enable":                                       "true",
		"traefik.http.routers.myapp.rule":                      "Host(`myapp.test`)",
		"traefik.http.services.myapp.loadbalancer.server.port": "80",
	}

	for k, v := range expected {
		got, ok := labels[k]
		if !ok {
			t.Errorf("missing label %q", k)
		} else if got != v {
			t.Errorf("label %q = %q, want %q", k, got, v)
		}
	}
}
