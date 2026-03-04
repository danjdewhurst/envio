package proxy

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// CertsDir returns the directory where mkcert certificates are stored.
func CertsDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".envio", "proxy", "certs")
}

// TLSConfigPath returns the path to the Traefik TLS dynamic config file.
func TLSConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".envio", "proxy", "traefik-tls.yml")
}

// IsMkcertInstalled checks whether mkcert is available on PATH.
func IsMkcertInstalled() bool {
	_, err := exec.LookPath("mkcert")
	return err == nil
}

// SetupTLS checks mkcert is installed and runs mkcert -install.
func SetupTLS() error {
	if !IsMkcertInstalled() {
		return fmt.Errorf("mkcert is not installed — install it first: https://github.com/FiloSottile/mkcert#installation")
	}

	cmd := exec.Command("mkcert", "-install")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// CertExists checks whether both cert and key files exist for the given domain.
func CertExists(domain string) bool {
	dir := CertsDir()
	cert := filepath.Join(dir, domain+".pem")
	key := filepath.Join(dir, domain+"-key.pem")
	_, errCert := os.Stat(cert)
	_, errKey := os.Stat(key)
	return errCert == nil && errKey == nil
}

// GenerateCert generates a TLS certificate for the given domain using mkcert.
func GenerateCert(domain string) error {
	if CertExists(domain) {
		return nil
	}

	dir := CertsDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create certs directory: %w", err)
	}

	cert := filepath.Join(dir, domain+".pem")
	key := filepath.Join(dir, domain+"-key.pem")

	cmd := exec.Command("mkcert", "-cert-file", cert, "-key-file", key, domain+".test")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// tlsConfig represents the Traefik dynamic TLS configuration file.
type tlsConfig struct {
	TLS tlsSection `yaml:"tls"`
}

type tlsSection struct {
	Certificates []tlsCertificate `yaml:"certificates"`
}

type tlsCertificate struct {
	CertFile string `yaml:"certFile"`
	KeyFile  string `yaml:"keyFile"`
}

// EnsureTLSConfig creates an empty traefik-tls.yml if it doesn't exist.
func EnsureTLSConfig() error {
	path := TLSConfigPath()
	if _, err := os.Stat(path); err == nil {
		return nil
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create proxy directory: %w", err)
	}

	cfg := tlsConfig{TLS: tlsSection{Certificates: []tlsCertificate{}}}
	return writeTLSConfig(path, &cfg)
}

// RegisterDomainCert adds a domain's certificate to traefik-tls.yml.
// Uses container-relative paths (/certs/...) so Traefik can find them.
func RegisterDomainCert(domain string) error {
	path := TLSConfigPath()

	cfg, err := readTLSConfig(path)
	if err != nil {
		cfg = &tlsConfig{TLS: tlsSection{Certificates: []tlsCertificate{}}}
	}

	certFile := fmt.Sprintf("/certs/%s.pem", domain)
	keyFile := fmt.Sprintf("/certs/%s-key.pem", domain)

	// Check if already registered
	for _, c := range cfg.TLS.Certificates {
		if c.CertFile == certFile {
			return nil
		}
	}

	cfg.TLS.Certificates = append(cfg.TLS.Certificates, tlsCertificate{
		CertFile: certFile,
		KeyFile:  keyFile,
	})

	return writeTLSConfig(path, cfg)
}

func readTLSConfig(path string) (*tlsConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg tlsConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func writeTLSConfig(path string, cfg *tlsConfig) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal TLS config: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}
