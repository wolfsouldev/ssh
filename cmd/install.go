package cmd

import (
	"github.com/astro/sshh/internal/installer"
	"github.com/astro/sshh/internal/ui"
	"github.com/spf13/cobra"
)

var installCommand = &cobra.Command{
	Use:   "install",
	Short: "Install sshh globally (adds to PATH)",
	Run: func(cmd *cobra.Command, args []string) {
		ui.PrintHeader("Installing SSHH")

		if installer.IsInstalled() {
			if !ui.Confirm("  SSHH is already installed. Reinstall?") {
				return
			}
		}

		if err := installer.Install(); err != nil {
			ui.PrintError("Installation failed: %v", err)
			return
		}

		ui.PrintSuccess("SSHH installed successfully!")
		ui.PrintInfo("Restart your terminal and type 'sshh' to get started.")
	},
}

var uninstallCommand = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall sshh from PATH",
	Run: func(cmd *cobra.Command, args []string) {
		ui.PrintHeader("Uninstalling SSHH")

		if !installer.IsInstalled() {
			ui.PrintError("SSHH is not installed")
			return
		}

		if !ui.Confirm("  Are you sure you want to uninstall SSHH?") {
			return
		}

		if err := installer.Uninstall(); err != nil {
			ui.PrintError("Uninstall failed: %v", err)
			return
		}

		ui.PrintSuccess("SSHH uninstalled")
		ui.PrintInfo("Your vault data in %%APPDATA%%\\sshh is preserved.")
	},
}
