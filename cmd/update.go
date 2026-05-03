package cmd

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/wolfsouldev/ssh/internal/ui"
	"github.com/wolfsouldev/ssh/internal/version"
	"github.com/spf13/cobra"
)

const (
	githubRepo = "wolfsouldev/ssh"
	githubAPI  = "https://api.github.com/repos/" + githubRepo + "/releases/latest"
)

type githubRelease struct {
	TagName string        `json:"tag_name"`
	Assets  []githubAsset `json:"assets"`
	HTMLURL string        `json:"html_url"`
}

type githubAsset struct {
	Name               string `json:"name"`
	Size               int64  `json:"size"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

var checkOnly bool

var updateCommand = &cobra.Command{
	Use:   "update",
	Short: "Update sshh to the latest version",
	Long:  "Check for updates and self-update from GitHub releases.",
	Run:   runUpdate,
}

func init() {
	updateCommand.Flags().BoolVar(&checkOnly, "check", false, "Only check for updates, don't install")
}

func runUpdate(cmd *cobra.Command, args []string) {
	ui.PrintBanner()
	ui.PrintHeader("Update SSHH")

	ui.PrintKeyValue("Current:", version.Version)
	ui.PrintInfo("Checking for updates...")
	fmt.Println()

	release, err := fetchLatestRelease()
	if err != nil {
		ui.PrintError("Failed to check for updates: %v", err)
		ui.PrintInfo("Check manually: https://github.com/%s/releases", githubRepo)
		return
	}

	latest := strings.TrimPrefix(release.TagName, "v")
	ui.PrintKeyValue("Latest:", latest)
	fmt.Println()

	if latest == version.Version {
		ui.PrintSuccess("You're already on the latest version!")
		return
	}

	if version.Version == "dev" {
		ui.PrintWarn("You're running a dev build.")
		ui.PrintInfo("Latest release: %s%sv%s%s", ui.Bold, ui.BrightCyan, latest, ui.Reset)
		ui.PrintInfo("Download: %s", release.HTMLURL)
		fmt.Println()
	}

	if checkOnly {
		ui.PrintInfo("Update available: %s → %s", version.Version, latest)
		ui.PrintInfo("Run 'sshh update' to install it.")
		return
	}

	if !ui.Confirm(fmt.Sprintf("  Update %s%s%s → %s%s%s?",
		ui.Dim, version.Version, ui.Reset,
		ui.Bold+ui.BrightGreen, latest, ui.Reset)) {
		ui.PrintInfo("Cancelled")
		return
	}

	// Find the right asset for this OS/arch
	assetName := buildAssetName(latest)
	var downloadURL string
	var assetSize int64
	for _, a := range release.Assets {
		if a.Name == assetName {
			downloadURL = a.BrowserDownloadURL
			assetSize = a.Size
			break
		}
	}

	if downloadURL == "" {
		ui.PrintError("No binary found for %s/%s", runtime.GOOS, runtime.GOARCH)
		ui.PrintInfo("Available assets:")
		for _, a := range release.Assets {
			fmt.Printf("    %s%s%s\n", ui.Dim, a.Name, ui.Reset)
		}
		ui.PrintInfo("Download manually: %s", release.HTMLURL)
		return
	}

	ui.PrintInfo("Downloading %s (%s)...", assetName, formatBytes(assetSize))

	tmpDir, err := os.MkdirTemp("", "sshh-update-*")
	if err != nil {
		ui.PrintError("Failed to create temp dir: %v", err)
		return
	}
	defer os.RemoveAll(tmpDir)

	archivePath := filepath.Join(tmpDir, assetName)
	if err := downloadFile(archivePath, downloadURL); err != nil {
		ui.PrintError("Download failed: %v", err)
		return
	}
	ui.PrintSuccess("Download complete")

	ui.PrintInfo("Extracting binary...")
	binaryPath, err := extractBinary(archivePath, tmpDir)
	if err != nil {
		ui.PrintError("Extract failed: %v", err)
		return
	}

	// Find current executable path
	currentBinary, err := os.Executable()
	if err != nil {
		ui.PrintError("Cannot find current binary: %v", err)
		return
	}
	currentBinary, err = filepath.EvalSymlinks(currentBinary)
	if err != nil {
		ui.PrintError("Cannot resolve binary path: %v", err)
		return
	}

	ui.PrintInfo("Installing to %s...", currentBinary)
	if err := replaceBinary(currentBinary, binaryPath); err != nil {
		ui.PrintError("Update failed: %v", err)
		ui.PrintWarn("You may need to run with sudo on Linux:")
		fmt.Printf("    %ssudo sshh update%s\n", ui.BrightWhite, ui.Reset)
		ui.PrintInfo("Or download manually: %s", release.HTMLURL)
		return
	}

	fmt.Println()
	ui.PrintSuccess("Updated to v%s!", latest)
	ui.PrintInfo("Restart sshh to use the new version.")
	fmt.Println()
}

// ── GitHub API ──

func fetchLatestRelease() (*githubRelease, error) {
	req, err := http.NewRequest("GET", githubAPI, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "sshh/"+version.Version)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("network error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("no releases found")
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("GitHub API returned %d", resp.StatusCode)
	}

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return &release, nil
}

// ── Download & Extract ──

func buildAssetName(ver string) string {
	ext := "tar.gz"
	if runtime.GOOS == "windows" {
		ext = "zip"
	}
	return fmt.Sprintf("sshh_%s_%s_%s.%s", ver, runtime.GOOS, runtime.GOARCH, ext)
}

func downloadFile(destPath, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("server returned %d", resp.StatusCode)
	}

	f, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	return err
}

func extractBinary(archivePath, destDir string) (string, error) {
	binaryName := "sshh"
	if runtime.GOOS == "windows" {
		binaryName = "sshh.exe"
	}

	if strings.HasSuffix(archivePath, ".zip") {
		return extractFromZip(archivePath, destDir, binaryName)
	}
	return extractFromTarGz(archivePath, destDir, binaryName)
}

func extractFromTarGz(archivePath, destDir, binaryName string) (string, error) {
	f, err := os.Open(archivePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return "", err
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}

		if filepath.Base(hdr.Name) == binaryName && hdr.Typeflag == tar.TypeReg {
			outPath := filepath.Join(destDir, binaryName)
			out, err := os.OpenFile(outPath, os.O_CREATE|os.O_WRONLY, 0755)
			if err != nil {
				return "", err
			}
			if _, err := io.Copy(out, tr); err != nil {
				out.Close()
				return "", err
			}
			out.Close()
			return outPath, nil
		}
	}

	return "", fmt.Errorf("binary '%s' not found in archive", binaryName)
}

func extractFromZip(archivePath, destDir, binaryName string) (string, error) {
	r, err := zip.OpenReader(archivePath)
	if err != nil {
		return "", err
	}
	defer r.Close()

	for _, f := range r.File {
		if filepath.Base(f.Name) == binaryName {
			rc, err := f.Open()
			if err != nil {
				return "", err
			}

			outPath := filepath.Join(destDir, binaryName)
			out, err := os.OpenFile(outPath, os.O_CREATE|os.O_WRONLY, 0755)
			if err != nil {
				rc.Close()
				return "", err
			}

			_, cpErr := io.Copy(out, rc)
			rc.Close()
			out.Close()
			if cpErr != nil {
				return "", cpErr
			}
			return outPath, nil
		}
	}

	return "", fmt.Errorf("binary '%s' not found in archive", binaryName)
}

// ── Binary Replacement ──

func replaceBinary(currentPath, newPath string) error {
	if runtime.GOOS == "windows" {
		// Windows can't overwrite a running exe, but can rename it
		oldPath := currentPath + ".old"
		os.Remove(oldPath)
		if err := os.Rename(currentPath, oldPath); err != nil {
			return fmt.Errorf("cannot rename current binary: %w", err)
		}
		if err := copyFile(newPath, currentPath); err != nil {
			// Try to restore
			os.Rename(oldPath, currentPath)
			return err
		}
		os.Remove(oldPath)
		return nil
	}

	// Linux/macOS: read new binary and overwrite
	newData, err := os.ReadFile(newPath)
	if err != nil {
		return err
	}
	return os.WriteFile(currentPath, newData, 0755)
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0755)
}

// ── Helpers ──

func formatBytes(b int64) string {
	switch {
	case b >= 1024*1024:
		return fmt.Sprintf("%.1f MB", float64(b)/(1024*1024))
	case b >= 1024:
		return fmt.Sprintf("%.1f KB", float64(b)/1024)
	default:
		return fmt.Sprintf("%d B", b)
	}
}
