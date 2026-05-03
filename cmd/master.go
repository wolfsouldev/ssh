package cmd

import (
	"fmt"

	"github.com/wolfsouldev/ssh/internal/ui"
	"github.com/wolfsouldev/ssh/internal/vault"
	"github.com/spf13/cobra"
)

var masterCommand = &cobra.Command{
	Use:   "master",
	Short: "Change the master password",
	Run: func(cmd *cobra.Command, args []string) {
		ui.PrintBanner()

		v, err := vault.New()
		if err != nil {
			ui.PrintError("Vault error: %v", err)
			return
		}

		if !v.Exists() {
			ui.PrintError("No vault found. Use 'sshh add' to create one.")
			return
		}

		ui.PrintHeader("Change Master Password")

		oldPass, err := ui.ReadPassword(ui.PasswordPrompt("Current master password"))
		if err != nil {
			ui.PrintError("Error: %v", err)
			return
		}

		newPass, err := ui.ReadPassword(ui.PasswordPrompt("New master password"))
		if err != nil {
			ui.PrintError("Error: %v", err)
			return
		}

		confirmPass, err := ui.ReadPassword(ui.PasswordPrompt("Confirm new master password"))
		if err != nil {
			ui.PrintError("Error: %v", err)
			return
		}

		if string(newPass) != string(confirmPass) {
			ui.PrintError("New passwords don't match")
			return
		}

		if len(newPass) < 4 {
			ui.PrintError("Master password must be at least 4 characters")
			return
		}

		if err := v.ChangeMasterPassword(oldPass, newPass); err != nil {
			ui.PrintError("%v", err)
			return
		}

		fmt.Println()
		ui.PrintSuccess("Master password changed successfully")
	},
}
