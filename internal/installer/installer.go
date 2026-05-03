package installer

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// InstallDir returns the installation directory.
func InstallDir() string {
	appData := os.Getenv("LOCALAPPDATA")
	if appData == "" {
		home, _ := os.UserHomeDir()
		appData = filepath.Join(home, "AppData", "Local")
	}
	return filepath.Join(appData, "sshh", "bin")
}

// Install copies the current executable to the install directory and adds it to PATH.
func Install() error {
	installDir := InstallDir()

	// Create install directory
	if err := os.MkdirAll(installDir, 0755); err != nil {
		return fmt.Errorf("creating install directory: %w", err)
	}

	// Get current executable path
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("getting executable path: %w", err)
	}

	destPath := filepath.Join(installDir, "sshh.exe")

	// Copy executable
	input, err := os.ReadFile(exePath)
	if err != nil {
		return fmt.Errorf("reading executable: %w", err)
	}

	if err := os.WriteFile(destPath, input, 0755); err != nil {
		return fmt.Errorf("writing executable: %w", err)
	}

	fmt.Printf("  Installed to: %s\n", destPath)

	// Add to PATH if not already there
	if err := addToPath(installDir); err != nil {
		return fmt.Errorf("adding to PATH: %w", err)
	}

	return nil
}

// Uninstall removes the installed executable and cleans PATH.
func Uninstall() error {
	installDir := InstallDir()
	destPath := filepath.Join(installDir, "sshh.exe")

	if _, err := os.Stat(destPath); os.IsNotExist(err) {
		return fmt.Errorf("sshh is not installed at %s", destPath)
	}

	if err := os.Remove(destPath); err != nil {
		return fmt.Errorf("removing executable: %w", err)
	}

	fmt.Printf("  Removed: %s\n", destPath)

	if err := removeFromPath(installDir); err != nil {
		fmt.Printf("  Warning: could not remove from PATH: %v\n", err)
	}

	return nil
}

// IsInstalled checks if sshh is installed in the install directory.
func IsInstalled() bool {
	destPath := filepath.Join(InstallDir(), "sshh.exe")
	_, err := os.Stat(destPath)
	return err == nil
}

func addToPath(dir string) error {
	// Read current user PATH from registry
	out, err := exec.Command("reg", "query", "HKCU\\Environment", "/v", "Path").Output()
	if err != nil {
		// PATH doesn't exist yet, create it
		cmd := exec.Command("reg", "add", "HKCU\\Environment", "/v", "Path", "/t", "REG_EXPAND_SZ", "/d", dir, "/f")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("creating PATH entry: %w", err)
		}
		broadcastChange()
		fmt.Println("  Added to user PATH")
		return nil
	}

	currentPath := parseRegOutput(string(out))
	if strings.Contains(strings.ToLower(currentPath), strings.ToLower(dir)) {
		fmt.Println("  Already in PATH")
		return nil
	}

	newPath := currentPath + ";" + dir
	cmd := exec.Command("reg", "add", "HKCU\\Environment", "/v", "Path", "/t", "REG_EXPAND_SZ", "/d", newPath, "/f")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("updating PATH: %w", err)
	}

	broadcastChange()
	fmt.Println("  Added to user PATH (restart terminal to take effect)")
	return nil
}

func removeFromPath(dir string) error {
	out, err := exec.Command("reg", "query", "HKCU\\Environment", "/v", "Path").Output()
	if err != nil {
		return nil // No PATH to clean
	}

	currentPath := parseRegOutput(string(out))
	parts := strings.Split(currentPath, ";")
	filtered := make([]string, 0, len(parts))
	for _, p := range parts {
		if !strings.EqualFold(strings.TrimSpace(p), dir) {
			filtered = append(filtered, p)
		}
	}

	newPath := strings.Join(filtered, ";")
	cmd := exec.Command("reg", "add", "HKCU\\Environment", "/v", "Path", "/t", "REG_EXPAND_SZ", "/d", newPath, "/f")
	return cmd.Run()
}

func parseRegOutput(output string) string {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "REG_EXPAND_SZ") || strings.Contains(line, "REG_SZ") {
			parts := strings.SplitN(line, "REG_EXPAND_SZ", 2)
			if len(parts) < 2 {
				parts = strings.SplitN(line, "REG_SZ", 2)
			}
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1])
			}
		}
	}
	return ""
}

func broadcastChange() {
	// Notify Windows that environment variables have changed
	exec.Command("cmd", "/c", "setx", "SSHH_INSTALLED", "1").Run()
}
