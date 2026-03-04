package proxy

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

const hostsMarkerBegin = "# BEGIN envio"
const hostsMarkerEnd = "# END envio"

// IsDNSConfigured checks if /etc/hosts contains envio-managed entries.
func IsDNSConfigured() bool {
	data, err := os.ReadFile("/etc/hosts")
	if err != nil {
		return false
	}
	return strings.Contains(string(data), hostsMarkerBegin)
}

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

// SetupDNS prints instructions for adding .test domains to /etc/hosts.
// On macOS, the .test TLD is blocked by the system resolver, so /etc/hosts
// is the only reliable way to resolve .test domains.
func SetupDNS() error {
	if runtime.GOOS != "darwin" {
		fmt.Println("Add entries to /etc/hosts for each .test domain:")
		fmt.Println("  127.0.0.1  myapp.test")
		return nil
	}

	fmt.Println("On macOS, .test domains must be added to /etc/hosts.")
	fmt.Println("The 'envio init' command manages this automatically.")
	fmt.Println()

	// Show current envio entries if any
	data, err := os.ReadFile("/etc/hosts")
	if err != nil {
		return fmt.Errorf("failed to read /etc/hosts: %w", err)
	}

	if strings.Contains(string(data), hostsMarkerBegin) {
		fmt.Println("Current envio entries in /etc/hosts:")
		inBlock := false
		for _, line := range strings.Split(string(data), "\n") {
			if line == hostsMarkerBegin {
				inBlock = true
				continue
			}
			if line == hostsMarkerEnd {
				inBlock = false
				continue
			}
			if inBlock {
				fmt.Println("  " + line)
			}
		}
	} else {
		fmt.Println("No envio entries found in /etc/hosts yet.")
		fmt.Println("Run 'envio init' to create a project with a .test domain.")
	}

	return nil
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
