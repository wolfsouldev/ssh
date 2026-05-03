package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/wolfsouldev/ssh/internal/crypto"
	"github.com/wolfsouldev/ssh/internal/ui"
	"github.com/wolfsouldev/ssh/internal/vault"
	"github.com/spf13/cobra"
)

var exportCommand = &cobra.Command{
	Use:   "export <file>",
	Short: "Export vault to a file (encrypted)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		outFile := args[0]

		v, err := vault.New()
		if err != nil {
			ui.PrintError("Vault error: %v", err)
			return
		}

		if !v.Exists() {
			ui.PrintError("No vault found")
			return
		}

		masterPass, err := ui.ReadPassword("  🔑 Master password: ")
		if err != nil {
			ui.PrintError("Error: %v", err)
			return
		}

		// Verify master password by loading
		data, err := v.Load(masterPass)
		if err != nil {
			ui.PrintError("%v", err)
			return
		}

		// Re-encrypt with same password for export
		jsonData, err := marshalVaultData(data)
		if err != nil {
			ui.PrintError("Error: %v", err)
			return
		}

		encData, err := encryptForExport(jsonData, masterPass)
		if err != nil {
			ui.PrintError("Error: %v", err)
			return
		}

		if err := os.WriteFile(outFile, encData, 0600); err != nil {
			ui.PrintError("Error writing file: %v", err)
			return
		}

		ui.PrintSuccess("Vault exported to %s (%d credentials)", outFile, len(data.Credentials))
	},
}

var importCommand = &cobra.Command{
	Use:   "import <file>",
	Short: "Import vault from an exported file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		inFile := args[0]

		v, err := vault.New()
		if err != nil {
			ui.PrintError("Vault error: %v", err)
			return
		}

		if v.Exists() {
			if !ui.Confirm("  A vault already exists. Import will REPLACE it. Continue?") {
				return
			}
		}

		encData, err := os.ReadFile(inFile)
		if err != nil {
			ui.PrintError("Error reading file: %v", err)
			return
		}

		importPass, err := ui.ReadPassword("  Password for imported file: ")
		if err != nil {
			ui.PrintError("Error: %v", err)
			return
		}

		plaintext, err := decryptForImport(encData, importPass)
		if err != nil {
			ui.PrintError("%v", err)
			return
		}

		data, err := unmarshalVaultData(plaintext)
		if err != nil {
			ui.PrintError("Invalid vault data: %v", err)
			return
		}

		// Ask for master password for the new local vault
		var masterPass []byte
		if v.Exists() {
			masterPass, err = ui.ReadPassword("  🔑 Your local master password: ")
			if err != nil {
				ui.PrintError("Error: %v", err)
				return
			}
			// Verify
			if _, err := v.Load(masterPass); err != nil {
				ui.PrintError("%v", err)
				return
			}
		} else {
			masterPass, err = ui.ReadPassword("  Set a master password for local vault: ")
			if err != nil {
				ui.PrintError("Error: %v", err)
				return
			}
			confirm, err := ui.ReadPassword("  Confirm master password: ")
			if err != nil {
				ui.PrintError("Error: %v", err)
				return
			}
			if string(masterPass) != string(confirm) {
				ui.PrintError("Passwords don't match")
				return
			}
		}

		// Save as new vault
		if err := saveImportedVault(v, data, masterPass); err != nil {
			ui.PrintError("Error: %v", err)
			return
		}

		ui.PrintSuccess("Imported %d credentials", len(data.Credentials))
	},
}

func marshalVaultData(data *vault.VaultData) ([]byte, error) {
	return json.MarshalIndent(data, "", "  ")
}

func encryptForExport(data []byte, pass []byte) ([]byte, error) {
	return crypto.Encrypt(data, pass)
}

func decryptForImport(data []byte, pass []byte) ([]byte, error) {
	return crypto.Decrypt(data, pass)
}

func unmarshalVaultData(data []byte) (*vault.VaultData, error) {
	var vd vault.VaultData
	if err := json.Unmarshal(data, &vd); err != nil {
		return nil, err
	}
	return &vd, nil
}

func saveImportedVault(v *vault.Vault, data *vault.VaultData, masterPass []byte) error {
	if err := v.Init(masterPass); err != nil {
		return fmt.Errorf("creating vault: %w", err)
	}
	for _, cred := range data.Credentials {
		if err := v.AddCredential(masterPass, cred); err != nil {
			return fmt.Errorf("adding credential '%s': %w", cred.Alias, err)
		}
	}
	return nil
}
