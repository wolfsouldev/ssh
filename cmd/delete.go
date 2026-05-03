package cmd

import (
	"fmt"

	"github.com/wolfsouldev/ssh/internal/ui"
	"github.com/wolfsouldev/ssh/internal/vault"
	"github.com/spf13/cobra"
)

var deleteCommand = &cobra.Command{
	Use:   "delete <alias>",
	Short: "Delete a saved SSH credential",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		alias := args[0]
		ui.PrintBanner()

		v, err := vault.New()
		if err != nil {
			ui.PrintError("Vault error: %v", err)
			return
		}

		if !v.Exists() {
			ui.PrintError("No vault found")
			return
		}

		masterPass, err := ui.ReadPassword(ui.PasswordPrompt("Master password"))
		if err != nil {
			ui.PrintError("Error: %v", err)
			return
		}

		// Verify the credential exists
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

		ui.PrintHeader("Delete Credential")
		ui.PrintKeyValue("Alias:", cred.Alias)
		ui.PrintKeyValue("Target:", cred.User+"@"+cred.Host)
		fmt.Println()

		if !ui.Confirm("  " + ui.Red + "Delete this credential?" + ui.Reset) {
			ui.PrintInfo("Cancelled")
			return
		}

		if err := v.DeleteCredential(masterPass, alias); err != nil {
			ui.PrintError("Error: %v", err)
			return
		}

		fmt.Println()
		ui.PrintSuccess("Credential '%s' deleted", alias)
	},
}
