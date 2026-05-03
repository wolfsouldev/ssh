package cmd

import (
	"fmt"

	"github.com/astro/sshh/internal/ui"
	"github.com/astro/sshh/internal/vault"
	"github.com/spf13/cobra"
)

var listCommand = &cobra.Command{
	Use:   "list",
	Short: "List all saved SSH credentials",
	Run: func(cmd *cobra.Command, args []string) {
		v, err := vault.New()
		if err != nil {
			ui.PrintError("Vault error: %v", err)
			return
		}

		if !v.Exists() {
			ui.PrintInfo("No vault found. Use 'sshh add' to create one.")
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

		ui.PrintHeader("Saved Credentials")

		headers := []string{"ALIAS", "USER", "HOST", "PORT", "TYPE", "UPDATED"}
		rows := make([][]string, 0, len(data.Credentials))
		for _, c := range data.Credentials {
			rows = append(rows, []string{
				c.Alias,
				c.User,
				c.Host,
				fmt.Sprintf("%d", c.Port),
				string(c.AuthType),
				c.UpdatedAt[:10],
			})
		}

		ui.PrintTable(headers, rows)
		fmt.Printf("\n  Total: %d credential(s)\n\n", len(data.Credentials))
	},
}
