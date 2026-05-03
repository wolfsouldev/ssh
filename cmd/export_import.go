package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/wolfsouldev/ssh/internal/crypto"
	"github.com/wolfsouldev/ssh/internal/ui"
	"github.com/wolfsouldev/ssh/internal/vault"
	"github.com/wolfsouldev/ssh/internal/version"
	"github.com/spf13/cobra"
)

// ── Export File Format ──
// [8 bytes magic: "SSHHEXP1"]
// [rest: AES-256-GCM encrypted JSON of ExportData]
// The export password is separate from the master password.

const exportMagic = "SSHHEXP1"

type ExportData struct {
	Version    int                `json:"version"`
	AppVersion string             `json:"app_version"`
	ExportedAt string             `json:"exported_at"`
	Count      int                `json:"count"`
	Creds      []vault.Credential `json:"credentials"`
}

// ── Export Command ──

var exportCommand = &cobra.Command{
	Use:   "export <file>",
	Short: "Export vault to a file (encrypted)",
	Long:  "Export all credentials to an encrypted .sshh file.\nUses a separate export password — your master password is never shared.",
	Args:  cobra.ExactArgs(1),
	Run:   runExport,
}

func runExport(cmd *cobra.Command, args []string) {
	outFile := args[0]
	if !strings.HasSuffix(outFile, ".sshh") {
		outFile += ".sshh"
	}

	ui.PrintBanner()
	ui.PrintHeader("Export Vault")

	v, err := vault.New()
	if err != nil {
		ui.PrintError("Vault error: %v", err)
		return
	}

	if !v.Exists() {
		ui.PrintError("No vault found. Nothing to export.")
		return
	}

	// Step 1: Unlock vault with master password
	masterPass, err := ui.ReadPassword(ui.PasswordPrompt("Master password"))
	if err != nil {
		ui.PrintError("Error: %v", err)
		return
	}

	data, err := v.Load(masterPass)
	if err != nil {
		ui.PrintError("%v", err)
		return
	}

	if len(data.Credentials) == 0 {
		ui.PrintWarn("Vault is empty. Nothing to export.")
		return
	}

	// Show what will be exported
	ui.PrintInfo("Credentials to export:")
	fmt.Println()
	headers := []string{"#", "NAME", "TARGET", "AUTH"}
	rows := make([][]string, 0, len(data.Credentials))
	for i, c := range data.Credentials {
		auth := "password"
		if c.AuthType == vault.AuthPrivateKey {
			auth = "key"
		}
		rows = append(rows, []string{
			fmt.Sprintf("%d", i+1),
			c.Alias,
			fmt.Sprintf("%s@%s:%d", c.User, c.Host, c.Port),
			auth,
		})
	}
	ui.PrintTable(headers, rows)
	fmt.Println()

	// Step 2: Set a separate export password
	ui.PrintInfo("Set a password to protect this export file.")
	ui.PrintWarn("Use a DIFFERENT password from your master password!")
	fmt.Println()

	exportPass, err := ui.ReadPassword(ui.PasswordPrompt("Export password"))
	if err != nil {
		ui.PrintError("Error: %v", err)
		return
	}

	if len(exportPass) < 4 {
		ui.PrintError("Export password must be at least 4 characters")
		return
	}

	confirmPass, err := ui.ReadPassword(ui.PasswordPrompt("Confirm export password"))
	if err != nil {
		ui.PrintError("Error: %v", err)
		return
	}

	if string(exportPass) != string(confirmPass) {
		ui.PrintError("Passwords don't match")
		return
	}

	// Step 3: Build export package
	pkg := ExportData{
		Version:    1,
		AppVersion: version.Version,
		ExportedAt: time.Now().Format(time.RFC3339),
		Count:      len(data.Credentials),
		Creds:      data.Credentials,
	}

	jsonData, err := json.Marshal(pkg)
	if err != nil {
		ui.PrintError("Error encoding data: %v", err)
		return
	}

	// Step 4: Encrypt with export password
	encData, err := crypto.Encrypt(jsonData, exportPass)
	if err != nil {
		ui.PrintError("Encryption failed: %v", err)
		return
	}

	// Step 5: Write magic header + encrypted data
	file, err := os.OpenFile(outFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		ui.PrintError("Error creating file: %v", err)
		return
	}
	defer file.Close()

	if _, err := file.Write([]byte(exportMagic)); err != nil {
		ui.PrintError("Error writing file: %v", err)
		return
	}
	if _, err := file.Write(encData); err != nil {
		ui.PrintError("Error writing file: %v", err)
		return
	}

	fmt.Println()
	ui.PrintSuccess("Exported %d credentials to %s", len(data.Credentials), outFile)
	ui.PrintKeyValue("File:", outFile)
	ui.PrintKeyValue("Encrypted:", "AES-256-GCM + Argon2id")
	ui.PrintInfo("Share this file + the export password to import on another machine.")
	fmt.Println()
}

// ── Import Command ──

var replaceFlag bool

var importCommand = &cobra.Command{
	Use:   "import <file>",
	Short: "Import vault from an exported file",
	Long:  "Import credentials from an encrypted .sshh export file.\nBy default, merges with existing vault. Use --replace to overwrite.",
	Args:  cobra.ExactArgs(1),
	Run:   runImport,
}

func init() {
	importCommand.Flags().BoolVar(&replaceFlag, "replace", false, "Replace existing vault instead of merging")
}

func runImport(cmd *cobra.Command, args []string) {
	inFile := args[0]

	ui.PrintBanner()
	ui.PrintHeader("Import Vault")

	// Step 1: Read and verify file format
	fileData, err := os.ReadFile(inFile)
	if err != nil {
		ui.PrintError("Error reading file: %v", err)
		return
	}

	if len(fileData) < len(exportMagic) || string(fileData[:len(exportMagic)]) != exportMagic {
		ui.PrintError("Not a valid SSHH export file")
		ui.PrintInfo("Expected .sshh file created by 'sshh export'")
		return
	}

	encData := fileData[len(exportMagic):]

	// Step 2: Decrypt with export password
	exportPass, err := ui.ReadPassword(ui.PasswordPrompt("Export password"))
	if err != nil {
		ui.PrintError("Error: %v", err)
		return
	}

	plaintext, err := crypto.Decrypt(encData, exportPass)
	if err != nil {
		ui.PrintError("Wrong export password or corrupted file")
		return
	}

	var pkg ExportData
	if err := json.Unmarshal(plaintext, &pkg); err != nil {
		ui.PrintError("Invalid export data: %v", err)
		return
	}

	// Step 3: Show import summary
	ui.PrintSuccess("File decrypted successfully!")
	fmt.Println()
	ui.PrintKeyValue("Exported by:", fmt.Sprintf("sshh v%s", pkg.AppVersion))
	ui.PrintKeyValue("Date:", pkg.ExportedAt)
	ui.PrintKeyValue("Credentials:", fmt.Sprintf("%d", pkg.Count))
	fmt.Println()

	headers := []string{"#", "NAME", "TARGET", "AUTH"}
	rows := make([][]string, 0, len(pkg.Creds))
	for i, c := range pkg.Creds {
		auth := "password"
		if c.AuthType == vault.AuthPrivateKey {
			auth = "key"
		}
		rows = append(rows, []string{
			fmt.Sprintf("%d", i+1),
			c.Alias,
			fmt.Sprintf("%s@%s:%d", c.User, c.Host, c.Port),
			auth,
		})
	}
	ui.PrintTable(headers, rows)
	fmt.Println()

	if !ui.Confirm("  Import these credentials?") {
		ui.PrintInfo("Cancelled")
		return
	}

	v, err := vault.New()
	if err != nil {
		ui.PrintError("Vault error: %v", err)
		return
	}

	if v.Exists() {
		importIntoExistingVault(v, pkg.Creds)
	} else {
		importIntoNewVault(v, pkg.Creds)
	}
}

// ── Import into existing vault (merge) ──

func importIntoExistingVault(v *vault.Vault, incoming []vault.Credential) {
	masterPass, err := ui.ReadPassword(ui.PasswordPrompt("Your local master password"))
	if err != nil {
		ui.PrintError("Error: %v", err)
		return
	}

	data, err := v.Load(masterPass)
	if err != nil {
		ui.PrintError("%v", err)
		return
	}

	if replaceFlag {
		ui.PrintWarn("Replace mode: all existing credentials will be removed.")
		if !ui.Confirm("  " + ui.Red + "Are you sure?" + ui.Reset) {
			ui.PrintInfo("Cancelled")
			return
		}

		if err := v.Init(masterPass); err != nil {
			ui.PrintError("Error resetting vault: %v", err)
			return
		}

		added := 0
		for _, cred := range incoming {
			if err := v.AddCredential(masterPass, cred); err != nil {
				ui.PrintError("Error adding '%s': %v", cred.Alias, err)
			} else {
				added++
			}
		}

		fmt.Println()
		ui.PrintSuccess("Replaced vault with %d credentials", added)
		return
	}

	// Merge mode (default)
	existingAliases := make(map[string]bool)
	for _, c := range data.Credentials {
		existingAliases[c.Alias] = true
	}

	added, skipped, overwritten := 0, 0, 0

	for _, cred := range incoming {
		if !existingAliases[cred.Alias] {
			// No conflict — add directly
			if err := v.AddCredential(masterPass, cred); err != nil {
				ui.PrintError("Error adding '%s': %v", cred.Alias, err)
			} else {
				added++
				ui.PrintSuccess("Added: %s", cred.Alias)
			}
			continue
		}

		// Conflict — ask user
		fmt.Println()
		ui.PrintWarn("Conflict: '%s' already exists locally", cred.Alias)
		ui.PrintKeyValue("  Local:", fmt.Sprintf("%s@%s:%d", cred.User, cred.Host, cred.Port))
		ui.PrintKeyValue("  Import:", fmt.Sprintf("%s@%s:%d", cred.User, cred.Host, cred.Port))
		fmt.Println()

		choice := ui.ReadLine(ui.InputPrompt("[s]kip / [o]verwrite / [r]ename"))
		choice = strings.ToLower(strings.TrimSpace(choice))

		switch choice {
		case "o", "overwrite":
			if err := v.DeleteCredential(masterPass, cred.Alias); err != nil {
				ui.PrintError("Error removing old '%s': %v", cred.Alias, err)
				continue
			}
			if err := v.AddCredential(masterPass, cred); err != nil {
				ui.PrintError("Error adding '%s': %v", cred.Alias, err)
			} else {
				overwritten++
				ui.PrintSuccess("Overwritten: %s", cred.Alias)
			}

		case "r", "rename":
			newAlias := ui.ReadLine(ui.InputPrompt("New name"))
			if newAlias == "" {
				ui.PrintWarn("Skipped '%s' (no name given)", cred.Alias)
				skipped++
				continue
			}
			cred.Alias = newAlias
			if err := v.AddCredential(masterPass, cred); err != nil {
				ui.PrintError("Error adding '%s': %v", newAlias, err)
			} else {
				added++
				ui.PrintSuccess("Added as: %s", newAlias)
			}

		default:
			skipped++
			ui.PrintInfo("Skipped: %s", cred.Alias)
		}
	}

	fmt.Println()
	ui.PrintDivider()
	ui.PrintKeyValue("Added:", fmt.Sprintf("%d", added))
	ui.PrintKeyValue("Overwritten:", fmt.Sprintf("%d", overwritten))
	ui.PrintKeyValue("Skipped:", fmt.Sprintf("%d", skipped))
	ui.PrintDivider()
	fmt.Println()
	ui.PrintSuccess("Import complete!")
}

// ── Import into new vault ──

func importIntoNewVault(v *vault.Vault, incoming []vault.Credential) {
	ui.PrintInfo("No vault found. Creating a new one.")
	fmt.Println()

	masterPass, err := ui.ReadPassword(ui.PasswordPrompt("Set a master password"))
	if err != nil {
		ui.PrintError("Error: %v", err)
		return
	}

	if len(masterPass) < 4 {
		ui.PrintError("Master password must be at least 4 characters")
		return
	}

	confirmPass, err := ui.ReadPassword(ui.PasswordPrompt("Confirm master password"))
	if err != nil {
		ui.PrintError("Error: %v", err)
		return
	}

	if string(masterPass) != string(confirmPass) {
		ui.PrintError("Passwords don't match")
		return
	}

	if err := v.Init(masterPass); err != nil {
		ui.PrintError("Error creating vault: %v", err)
		return
	}

	added := 0
	for _, cred := range incoming {
		if err := v.AddCredential(masterPass, cred); err != nil {
			ui.PrintError("Error adding '%s': %v", cred.Alias, err)
		} else {
			added++
		}
	}

	fmt.Println()
	ui.PrintSuccess("Vault created with %d credentials!", added)
}
