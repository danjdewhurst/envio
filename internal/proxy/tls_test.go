package proxy

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCertsDir(t *testing.T) {
	dir := CertsDir()
	if !strings.HasSuffix(dir, filepath.Join(".envio", "proxy", "certs")) {
		t.Errorf("expected path ending in .envio/proxy/certs, got %s", dir)
	}
}

func TestTLSConfigPath(t *testing.T) {
	path := TLSConfigPath()
	if !strings.HasSuffix(path, filepath.Join(".envio", "proxy", "traefik-tls.yml")) {
		t.Errorf("expected path ending in .envio/proxy/traefik-tls.yml, got %s", path)
	}
}

func TestCertExistsFalse(t *testing.T) {
	if CertExists("nonexistent-domain-xyz") {
		t.Error("CertExists should return false for nonexistent domain")
	}
}

func TestEnsureTLSConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "traefik-tls.yml")

	// Write to a custom path by calling writeTLSConfig directly
	cfg := &tlsConfig{TLS: tlsSection{Certificates: []tlsCertificate{}}}
	if err := writeTLSConfig(path, cfg); err != nil {
		t.Fatalf("writeTLSConfig failed: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read TLS config: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "tls:") {
		t.Error("TLS config should contain 'tls:' key")
	}
}

func TestRegisterDomainCert(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "traefik-tls.yml")

	// Start with empty config
	cfg := &tlsConfig{TLS: tlsSection{Certificates: []tlsCertificate{}}}
	if err := writeTLSConfig(path, cfg); err != nil {
		t.Fatalf("writeTLSConfig failed: %v", err)
	}

	// Register a domain
	registerDomain := func(domain string) {
		t.Helper()
		loaded, err := readTLSConfig(path)
		if err != nil {
			loaded = &tlsConfig{TLS: tlsSection{Certificates: []tlsCertificate{}}}
		}

		certFile := "/certs/" + domain + ".pem"
		keyFile := "/certs/" + domain + "-key.pem"

		for _, c := range loaded.TLS.Certificates {
			if c.CertFile == certFile {
				return
			}
		}

		loaded.TLS.Certificates = append(loaded.TLS.Certificates, tlsCertificate{
			CertFile: certFile,
			KeyFile:  keyFile,
		})

		if err := writeTLSConfig(path, loaded); err != nil {
			t.Fatalf("writeTLSConfig failed: %v", err)
		}
	}

	// Register first domain
	registerDomain("myapp")

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read TLS config: %v", err)
	}
	content := string(data)

	if !strings.Contains(content, "/certs/myapp.pem") {
		t.Error("TLS config should contain /certs/myapp.pem")
	}
	if !strings.Contains(content, "/certs/myapp-key.pem") {
		t.Error("TLS config should contain /certs/myapp-key.pem")
	}

	// Register same domain again (idempotency)
	registerDomain("myapp")

	data, err = os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read TLS config: %v", err)
	}
	count := strings.Count(string(data), "/certs/myapp.pem")
	if count != 1 {
		t.Errorf("expected 1 occurrence of /certs/myapp.pem, got %d", count)
	}

	// Register a second domain (preserves existing)
	registerDomain("otherapp")

	data, err = os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read TLS config: %v", err)
	}
	content = string(data)

	if !strings.Contains(content, "/certs/myapp.pem") {
		t.Error("TLS config should still contain /certs/myapp.pem after adding second domain")
	}
	if !strings.Contains(content, "/certs/otherapp.pem") {
		t.Error("TLS config should contain /certs/otherapp.pem")
	}
}
