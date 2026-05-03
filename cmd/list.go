package cmd

import (
	"fmt"

	"github.com/wolfsouldev/ssh/internal/ui"
	"github.com/wolfsouldev/ssh/internal/vault"
	"github.com/spf13/cobra"
)

var listCommand = &cobra.Command{
	Use:   "list",
	Short: "List all saved SSH credentials",
	Run: func(cmd *cobra.Command, args []string) {
		ui.PrintBanner()

		v, err := vault.New()
		if err != nil {
			ui.PrintError("Vault error: %v", err)
			return
		}

		if !v.Exists() {
			ui.PrintWarn("No vault found. Use 'sshh add' to create one.")
			return
		}

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

		ui.PrintHeader("Saved Credentials")

		headers := []string{"ALIAS", "USER", "HOST", "PORT", "TYPE", "UPDATED"}
		rows := make([][]string, 0, len(data.Credentials))
		for _, c := range data.Credentials {
			updated := c.UpdatedAt
			if len(updated) >= 10 {
				updated = updated[:10]
			}
			rows = append(rows, []string{
				c.Alias,
				c.User,
				c.Host,
				fmt.Sprintf("%d", c.Port),
				string(c.AuthType),
				updated,
			})
		}

		ui.PrintTable(headers, rows)
		fmt.Printf("\n  %s%sTotal: %d credential(s)%s\n\n", ui.Dim, ui.BrightGreen, len(data.Credentials), ui.Reset)
	},
}
