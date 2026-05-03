package cmd

import (
	"fmt"

	"github.com/wolfsouldev/ssh/internal/installer"
	"github.com/wolfsouldev/ssh/internal/ui"
	"github.com/spf13/cobra"
)

var installCommand = &cobra.Command{
	Use:   "install",
	Short: "Install sshh globally (adds to PATH)",
	Run: func(cmd *cobra.Command, args []string) {
		ui.PrintBanner()
		ui.PrintHeader("Installing SSHH")

		if installer.IsInstalled() {
			ui.PrintWarn("SSHH is already installed.")
			if !ui.Confirm("  Reinstall?") {
				return
			}
		}

		if err := installer.Install(); err != nil {
			ui.PrintError("Installation failed: %v", err)
			return
		}

		fmt.Println()
		ui.PrintSuccess("SSHH installed successfully!")
		ui.PrintInfo("Restart your terminal and type 'sshh' to get started.")
	},
}

var uninstallCommand = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall sshh from PATH",
	Run: func(cmd *cobra.Command, args []string) {
		ui.PrintBanner()
		ui.PrintHeader("Uninstalling SSHH")

		if !installer.IsInstalled() {
			ui.PrintError("SSHH is not installed")
			return
		}

		if !ui.Confirm("  " + ui.Red + "Are you sure you want to uninstall SSHH?" + ui.Reset) {
			return
		}

		if err := installer.Uninstall(); err != nil {
			ui.PrintError("Uninstall failed: %v", err)
			return
		}

		fmt.Println()
		ui.PrintSuccess("SSHH uninstalled")
		ui.PrintInfo("Your vault data is preserved.")
	},
}
