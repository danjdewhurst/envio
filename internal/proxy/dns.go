package proxy

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const hostsMarkerBegin = "# BEGIN envio"
const hostsMarkerEnd = "# END envio"

// sudoCommand creates an exec.Cmd that runs with sudo and has stdin/stdout/stderr
// attached to the current terminal so the user can enter their password.
func sudoCommand(name string, args ...string) *exec.Cmd {
	fullArgs := append([]string{name}, args...)
	cmd := exec.Command("sudo", fullArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

// AddHostsEntry adds a .test domain to /etc/hosts inside the envio block.
// Creates the block if it doesn't exist. Requires sudo.
func AddHostsEntry(domain string) error {
	fullDomain := domain + ".test"

	data, err := os.ReadFile("/etc/hosts")
	if err != nil {
		return fmt.Errorf("failed to read /etc/hosts: %w", err)
	}
	content := string(data)

	entry := fmt.Sprintf("127.0.0.1  %s", fullDomain)

	// Already present
	if strings.Contains(content, entry) {
		return nil
	}

	var newContent string
	if strings.Contains(content, hostsMarkerBegin) {
		// Add to existing block (before the END marker)
		newContent = strings.Replace(content, hostsMarkerEnd, entry+"\n"+hostsMarkerEnd, 1)
	} else {
		// Create new block at end of file
		block := fmt.Sprintf("\n%s\n%s\n%s\n", hostsMarkerBegin, entry, hostsMarkerEnd)
		newContent = strings.TrimRight(content, "\n") + "\n" + block
	}

	return writeHosts(newContent)
}

// RemoveHostsEntry removes a .test domain from the envio block in /etc/hosts.
func RemoveHostsEntry(domain string) error {
	fullDomain := domain + ".test"

	data, err := os.ReadFile("/etc/hosts")
	if err != nil {
		return fmt.Errorf("failed to read /etc/hosts: %w", err)
	}
	content := string(data)

	entry := fmt.Sprintf("127.0.0.1  %s", fullDomain)
	if !strings.Contains(content, entry) {
		return nil
	}

	newContent := strings.Replace(content, entry+"\n", "", 1)

	// Clean up empty block
	emptyBlock := fmt.Sprintf("%s\n%s", hostsMarkerBegin, hostsMarkerEnd)
	newContent = strings.Replace(newContent, "\n"+emptyBlock+"\n", "", 1)

	return writeHosts(newContent)
}

// writeHosts writes content to /etc/hosts via a temp file and sudo mv.
func writeHosts(content string) error {
	tmp, err := os.CreateTemp("", "envio-hosts-*")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmp.Name()

	if _, err := tmp.WriteString(content); err != nil {
		_ = tmp.Close()
		_ = os.Remove(tmpPath)
		return fmt.Errorf("failed to write temp file: %w", err)
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("failed to close temp file: %w", err)
	}

	// Copy permissions from original
	if err := os.Chmod(tmpPath, 0644); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("failed to set permissions: %w", err)
	}

	if err := sudoCommand("mv", tmpPath, "/etc/hosts").Run(); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("failed to update /etc/hosts (requires sudo): %w", err)
	}

	return nil
}
