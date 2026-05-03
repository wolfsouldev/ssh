package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/astro/sshh/internal/ui"
	"github.com/astro/sshh/internal/vault"
	"github.com/spf13/cobra"
)

var addCommand = &cobra.Command{
	Use:   "add",
	Short: "Add a new SSH credential interactively",
	Run: func(cmd *cobra.Command, args []string) {
		v, err := vault.New()
		if err != nil {
			ui.PrintError("Vault error: %v", err)
			return
		}

		var masterPass []byte

		if !v.Exists() {
			ui.PrintInfo("No vault found. Creating a new one...")
			masterPass, err = ui.ReadPassword("  New master password: ")
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
			if err := v.Init(masterPass); err != nil {
				ui.PrintError("Error creating vault: %v", err)
				return
			}
			ui.PrintSuccess("Vault created")
		} else {
			masterPass, err = ui.ReadPassword("  🔑 Master password: ")
			if err != nil {
				ui.PrintError("Error: %v", err)
				return
			}
			// Verify master password
			if _, err := v.Load(masterPass); err != nil {
				ui.PrintError("%v", err)
				return
			}
		}

		ui.PrintHeader("Add New Credential")

		host := ui.ReadLine("  Host: ")
		if host == "" {
			ui.PrintError("Host is required")
			return
		}

		user := ui.ReadLine("  User: ")
		if user == "" {
			ui.PrintError("User is required")
			return
		}

		portStr := ui.ReadLine("  Port [22]: ")
		port := 22
		if portStr != "" {
			if p, err := strconv.Atoi(portStr); err == nil {
				port = p
			}
		}

		alias := ui.ReadLine(fmt.Sprintf("  Alias [%s@%s]: ", user, host))
		if alias == "" {
			alias = fmt.Sprintf("%s@%s", user, host)
		}

		authType := ui.ReadLine("  Auth type (password/key) [password]: ")
		authType = strings.ToLower(strings.TrimSpace(authType))

		cred := vault.Credential{
			Alias: alias,
			Host:  host,
			Port:  port,
			User:  user,
		}

		if authType == "key" || authType == "privatekey" {
			cred.AuthType = vault.AuthPrivateKey

			keyPath := ui.ReadLine("  Path to private key: ")
			if keyPath == "" {
				ui.PrintError("Key path is required")
				return
			}

			keyData, err := os.ReadFile(keyPath)
			if err != nil {
				ui.PrintError("Error reading key file: %v", err)
				return
			}
			cred.PrivateKey = string(keyData)

			if ui.Confirm("  Does the key have a passphrase?") {
				keyPass, err := ui.ReadPassword("  Key passphrase: ")
				if err != nil {
					ui.PrintError("Error: %v", err)
					return
				}
				cred.KeyPass = string(keyPass)
			}
		} else {
			cred.AuthType = vault.AuthPassword
			pass, err := ui.ReadPassword("  SSH Password: ")
			if err != nil {
				ui.PrintError("Error: %v", err)
				return
			}
			cred.Password = string(pass)
		}

		if err := v.AddCredential(masterPass, cred); err != nil {
			ui.PrintError("Error: %v", err)
			return
		}

		ui.PrintSuccess("Credential '%s' saved successfully", alias)
	},
}
