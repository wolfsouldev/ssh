package cmd

import (
	"fmt"
	"runtime"

	"github.com/wolfsouldev/ssh/internal/ui"
	"github.com/wolfsouldev/ssh/internal/version"
	"github.com/spf13/cobra"
)

var versionCommand = &cobra.Command{
	Use:   "version",
	Short: "Show sshh version information",
	Run: func(cmd *cobra.Command, args []string) {
		ui.PrintBanner()
		ui.PrintHeader("Version Info")
		ui.PrintKeyValue("Version:", version.Version)
		ui.PrintKeyValue("Commit:", version.Commit)
		ui.PrintKeyValue("Built:", version.Date)
		ui.PrintKeyValue("Go:", runtime.Version())
		ui.PrintKeyValue("OS/Arch:", fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH))
		fmt.Println()
	},
}
