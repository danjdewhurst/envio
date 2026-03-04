package proxy

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

const networkName = "envio-proxy"

// NetworkName returns the Docker network name used by the proxy.
func NetworkName() string {
	return networkName
}

// ProxyDir returns the directory where the proxy compose file is stored.
func ProxyDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".envio", "proxy")
}

// WriteComposeFile writes the Traefik docker-compose.yml to the given directory.
func WriteComposeFile(dir string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create proxy directory: %w", err)
	}

	// Ensure traefik-tls.yml exists so Docker doesn't create a directory mount
	if err := EnsureTLSConfig(); err != nil {
		return fmt.Errorf("failed to ensure TLS config: %w", err)
	}

	certsDir := CertsDir()
	tlsConfigPath := TLSConfigPath()

	content := fmt.Sprintf(`services:
  traefik:
    image: traefik:v3.3
    command:
      - "--providers.docker=true"
      - "--providers.docker.exposedbydefault=false"
      - "--providers.docker.network=envio-proxy"
      - "--entrypoints.web.address=:80"
      - "--entrypoints.websecure.address=:443"
      - "--providers.file.filename=/etc/traefik/tls.yml"
      - "--providers.file.watch=true"
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - "/var/run/docker.sock:/var/run/docker.sock:ro"
      - "%s:/certs:ro"
      - "%s:/etc/traefik/tls.yml:ro"
    networks:
      - envio-proxy
    restart: unless-stopped
networks:
  envio-proxy:
    external: true
`, certsDir, tlsConfigPath)

	return os.WriteFile(filepath.Join(dir, "docker-compose.yml"), []byte(content), 0644)
}

// EnsureNetwork creates the envio-proxy Docker network if it doesn't exist.
func EnsureNetwork() error {
	cmd := exec.Command("docker", "network", "inspect", networkName)
	if cmd.Run() == nil {
		return nil // network already exists
	}

	cmd = exec.Command("docker", "network", "create", networkName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Start starts the Traefik proxy.
func Start() error {
	if err := EnsureNetwork(); err != nil {
		return fmt.Errorf("failed to create proxy network: %w", err)
	}

	dir := ProxyDir()
	if err := WriteComposeFile(dir); err != nil {
		return err
	}

	cmd := exec.Command("docker", "compose", "up", "-d")
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Stop stops the Traefik proxy.
func Stop() error {
	dir := ProxyDir()
	cmd := exec.Command("docker", "compose", "down")
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// IsRunning checks if the Traefik proxy is currently running.
func IsRunning() bool {
	dir := ProxyDir()
	composePath := filepath.Join(dir, "docker-compose.yml")
	if _, err := os.Stat(composePath); err != nil {
		return false
	}

	cmd := exec.Command("docker", "compose", "ps", "--status", "running", "-q")
	cmd.Dir = dir
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(output)) != ""
}

// SanitiseDomain converts a string into a valid domain label: lowercase,
// non-alphanumeric characters replaced with hyphens, collapsed and trimmed.
func SanitiseDomain(s string) string {
	s = strings.ToLower(s)
	re := regexp.MustCompile(`[^a-z0-9]+`)
	s = re.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	return s
}

// TraefikLabels returns the Docker labels needed to route traffic via Traefik.
func TraefikLabels(domain string, port int) map[string]string {
	return map[string]string{
		"traefik.enable": "true",
		// HTTP router — redirects to HTTPS
		fmt.Sprintf("traefik.http.routers.%s.rule", domain):        fmt.Sprintf("Host(`%s.test`)", domain),
		fmt.Sprintf("traefik.http.routers.%s.entrypoints", domain): "web",
		fmt.Sprintf("traefik.http.routers.%s.middlewares", domain): "redirect-to-https",
		// HTTPS router
		fmt.Sprintf("traefik.http.routers.%s-tls.rule", domain):        fmt.Sprintf("Host(`%s.test`)", domain),
		fmt.Sprintf("traefik.http.routers.%s-tls.entrypoints", domain): "websecure",
		fmt.Sprintf("traefik.http.routers.%s-tls.tls", domain):         "true",
		// Service
		fmt.Sprintf("traefik.http.services.%s.loadbalancer.server.port", domain): fmt.Sprintf("%d", port),
		// Redirect middleware
		"traefik.http.middlewares.redirect-to-https.redirectscheme.scheme":    "https",
		"traefik.http.middlewares.redirect-to-https.redirectscheme.permanent": "true",
	}
}
