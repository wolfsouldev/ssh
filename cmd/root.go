package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/wolfsouldev/ssh/internal/sshclient"
	"github.com/wolfsouldev/ssh/internal/ui"
	"github.com/wolfsouldev/ssh/internal/vault"
	"github.com/spf13/cobra"
)

var portFlag int

var rootCmd = &cobra.Command{
	Use:   "sshh [user@host]",
	Short: "SSHH - Secure SSH credential manager & connector",
	Long: `SSHH is a secure SSH credential manager that stores your SSH 
credentials encrypted with a master password. Connect to your servers 
instantly without remembering individual passwords.`,
	Args: cobra.MaximumNArgs(1),
	Run:  connectCmd,
}

func init() {
	rootCmd.Flags().IntVarP(&portFlag, "port", "p", 22, "SSH port")

	rootCmd.AddCommand(listCommand)
	rootCmd.AddCommand(addCommand)
	rootCmd.AddCommand(editCommand)
	rootCmd.AddCommand(deleteCommand)
	rootCmd.AddCommand(masterCommand)
	rootCmd.AddCommand(installCommand)
	rootCmd.AddCommand(uninstallCommand)
	rootCmd.AddCommand(exportCommand)
	rootCmd.AddCommand(importCommand)
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func connectCmd(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		cmd.Help()
		return
	}

	target := args[0]
	user, host, port := parseTarget(target)

	if portFlag != 22 {
		port = portFlag
	}

	if user == "" || host == "" {
		ui.PrintError("Invalid target. Use: sshh user@host")
		return
	}

	fmt.Printf("\n  Connecting to %s@%s:%d ...\n", user, host, port)

	v, err := vault.New()
	if err != nil {
		ui.PrintError("Vault error: %v", err)
		return
	}

	if v.Exists() {
		// Try to find saved credentials
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

		cred := v.FindByHostUser(data, host, user)
		if cred != nil {
			ui.PrintSuccess("Credentials found for %s@%s", user, host)
			if err := connectWithCred(cred); err != nil {
				ui.PrintError("Connection failed: %v", err)
			}
			return
		}

		// No saved credentials - ask for password and connect
		ui.PrintInfo("No saved credentials for %s@%s", user, host)
		connectAndOfferSave(v, masterPass, user, host, port)
	} else {
		// No vault exists yet
		ui.PrintInfo("No vault found. Connecting without saved credentials.")
		connectNewNoVault(v, user, host, port)
	}
}

func connectWithCred(cred *vault.Credential) error {
	switch cred.AuthType {
	case vault.AuthPrivateKey:
		ui.PrintInfo("Authenticating with private key...")
		return sshclient.ConnectWithKey(cred.Host, cred.Port, cred.User, cred.PrivateKey, cred.KeyPass)
	default:
		ui.PrintInfo("Authenticating with password...")
		return sshclient.ConnectWithPassword(cred.Host, cred.Port, cred.User, cred.Password)
	}
}

func connectAndOfferSave(v *vault.Vault, masterPass []byte, user, host string, port int) {
	password, err := ui.ReadPassword(fmt.Sprintf("  Password for %s@%s: ", user, host))
	if err != nil {
		ui.PrintError("Error: %v", err)
		return
	}

	err = sshclient.ConnectWithPassword(host, port, user, string(password))
	if err != nil {
		ui.PrintError("Connection failed: %v", err)
		return
	}

	// After disconnecting, offer to save
	fmt.Println()
	if ui.Confirm("  Save credentials for " + user + "@" + host + "?") {
		alias := ui.ReadLine("  Alias (leave empty for " + user + "@" + host + "): ")
		if alias == "" {
			alias = fmt.Sprintf("%s@%s", user, host)
		}

		cred := vault.Credential{
			Alias:    alias,
			Host:     host,
			Port:     port,
			User:     user,
			AuthType: vault.AuthPassword,
			Password: string(password),
		}

		if err := v.AddCredential(masterPass, cred); err != nil {
			ui.PrintError("Error saving: %v", err)
		} else {
			ui.PrintSuccess("Credentials saved as '%s'", alias)
		}
	}
}

func connectNewNoVault(v *vault.Vault, user, host string, port int) {
	password, err := ui.ReadPassword(fmt.Sprintf("  Password for %s@%s: ", user, host))
	if err != nil {
		ui.PrintError("Error: %v", err)
		return
	}

	err = sshclient.ConnectWithPassword(host, port, user, string(password))
	if err != nil {
		ui.PrintError("Connection failed: %v", err)
		return
	}

	// After disconnecting, offer to save
	fmt.Println()
	if ui.Confirm("  Save credentials?") {
		ui.PrintInfo("First, set up a master password to encrypt your vault.")
		masterPass, err := ui.ReadPassword("  New master password: ")
		if err != nil {
			ui.PrintError("Error: %v", err)
			return
		}
		confirmPass, err := ui.ReadPassword("  Confirm master password: ")
		if err != nil {
			ui.PrintError("Error: %v", err)
			return
		}
		if string(masterPass) != string(confirmPass) {
			ui.PrintError("Passwords don't match")
			return
		}

		if err := v.Init(masterPass); err != nil {
			ui.PrintError("Error creating vault: %v", err)
			return
		}

		alias := ui.ReadLine("  Alias (leave empty for " + user + "@" + host + "): ")
		if alias == "" {
			alias = fmt.Sprintf("%s@%s", user, host)
		}

		cred := vault.Credential{
			Alias:    alias,
			Host:     host,
			Port:     port,
			User:     user,
			AuthType: vault.AuthPassword,
			Password: string(password),
		}

		if err := v.AddCredential(masterPass, cred); err != nil {
			ui.PrintError("Error saving: %v", err)
		} else {
			ui.PrintSuccess("Vault created and credentials saved as '%s'", alias)
		}
	}
}

func parseTarget(target string) (user, host string, port int) {
	port = 22

	// Handle user@host:port or user@host
	if at := strings.Index(target, "@"); at > 0 {
		user = target[:at]
		rest := target[at+1:]

		if colon := strings.LastIndex(rest, ":"); colon > 0 {
			host = rest[:colon]
			if p, err := strconv.Atoi(rest[colon+1:]); err == nil {
				port = p
			}
		} else {
			host = rest
		}
	} else {
		host = target
	}

	return
}
