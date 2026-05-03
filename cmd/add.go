package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/wolfsouldev/ssh/internal/ui"
	"github.com/wolfsouldev/ssh/internal/vault"
	"github.com/spf13/cobra"
)

var addCommand = &cobra.Command{
	Use:   "add",
	Short: "Add a new SSH credential interactively",
	Run: func(cmd *cobra.Command, args []string) {
		ui.PrintBanner()

		v, err := vault.New()
		if err != nil {
			ui.PrintError("Vault error: %v", err)
			return
		}

		var masterPass []byte

		if !v.Exists() {
			ui.PrintWarn("No vault found. Creating a new one...")
			fmt.Println()
			masterPass, err = ui.ReadPassword(ui.PasswordPrompt("New master password"))
			if err != nil {
				ui.PrintError("Error: %v", err)
				return
			}
			confirm, err := ui.ReadPassword(ui.PasswordPrompt("Confirm master password"))
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
			masterPass, err = ui.ReadPassword(ui.PasswordPrompt("Master password"))
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

		host := ui.ReadLine(ui.InputPrompt("Host"))
		if host == "" {
			ui.PrintError("Host is required")
			return
		}

		user := ui.ReadLine(ui.InputPrompt("User"))
		if user == "" {
			ui.PrintError("User is required")
			return
		}

		portStr := ui.ReadLine(ui.InputPrompt("Port [22]"))
		port := 22
		if portStr != "" {
			if p, err := strconv.Atoi(portStr); err == nil {
				port = p
			}
		}

		alias := ui.ReadLine(ui.InputPrompt(fmt.Sprintf("Alias [%s@%s]", user, host)))
		if alias == "" {
			alias = fmt.Sprintf("%s@%s", user, host)
		}

		authType := ui.ReadLine(ui.InputPrompt("Auth type (password/key) [password]"))
		authType = strings.ToLower(strings.TrimSpace(authType))

		cred := vault.Credential{
			Alias: alias,
			Host:  host,
			Port:  port,
			User:  user,
		}

		if authType == "key" || authType == "privatekey" {
			cred.AuthType = vault.AuthPrivateKey

			keyPath := ui.ReadLine(ui.InputPrompt("Path to private key"))
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
				keyPass, err := ui.ReadPassword(ui.PasswordPrompt("Key passphrase"))
				if err != nil {
					ui.PrintError("Error: %v", err)
					return
				}
				cred.KeyPass = string(keyPass)
			}
		} else {
			cred.AuthType = vault.AuthPassword
			pass, err := ui.ReadPassword(ui.PasswordPrompt("SSH Password"))
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

		fmt.Println()
		ui.PrintSuccess("Credential '%s' saved successfully", alias)
	},
}
