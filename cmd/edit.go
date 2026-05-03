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

var editCommand = &cobra.Command{
	Use:   "edit <alias>",
	Short: "Edit an existing SSH credential",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		alias := args[0]

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

		data, err := v.Load(masterPass)
		if err != nil {
			ui.PrintError("%v", err)
			return
		}

		cred := v.FindByAlias(data, alias)
		if cred == nil {
			ui.PrintError("Credential '%s' not found", alias)
			return
		}

		ui.PrintHeader("Editing: " + alias)
		ui.PrintInfo("Press Enter to keep current value")

		err = v.UpdateCredential(masterPass, alias, func(c *vault.Credential) {
			newHost := ui.ReadLine(fmt.Sprintf("  Host [%s]: ", c.Host))
			if newHost != "" {
				c.Host = newHost
			}

			newUser := ui.ReadLine(fmt.Sprintf("  User [%s]: ", c.User))
			if newUser != "" {
				c.User = newUser
			}

			newPort := ui.ReadLine(fmt.Sprintf("  Port [%d]: ", c.Port))
			if newPort != "" {
				if p, err := strconv.Atoi(newPort); err == nil {
					c.Port = p
				}
			}

			newAlias := ui.ReadLine(fmt.Sprintf("  Alias [%s]: ", c.Alias))
			if newAlias != "" {
				c.Alias = newAlias
			}

			newAuthType := ui.ReadLine(fmt.Sprintf("  Auth type [%s]: ", c.AuthType))
			newAuthType = strings.ToLower(strings.TrimSpace(newAuthType))

			if newAuthType == "key" || newAuthType == "privatekey" {
				c.AuthType = vault.AuthPrivateKey
				keyPath := ui.ReadLine("  Path to new private key (empty to keep): ")
				if keyPath != "" {
					keyData, err := os.ReadFile(keyPath)
					if err != nil {
						ui.PrintError("Error reading key: %v", err)
					} else {
						c.PrivateKey = string(keyData)
					}
				}
				if ui.Confirm("  Update key passphrase?") {
					keyPass, err := ui.ReadPassword("  New key passphrase: ")
					if err == nil {
						c.KeyPass = string(keyPass)
					}
				}
			} else if newAuthType == "password" || (newAuthType == "" && c.AuthType == vault.AuthPassword) {
				if ui.Confirm("  Update password?") {
					pass, err := ui.ReadPassword("  New SSH password: ")
					if err == nil {
						c.Password = string(pass)
					}
				}
			}
		})

		if err != nil {
			ui.PrintError("Error: %v", err)
			return
		}

		ui.PrintSuccess("Credential updated")
	},
}
