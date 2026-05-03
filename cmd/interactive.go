package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/wolfsouldev/ssh/internal/ui"
	"github.com/wolfsouldev/ssh/internal/vault"
)

func runInteractive() {
	ui.PrintBanner()
	ui.PrintAstronaut()

	v, err := vault.New()
	if err != nil {
		ui.PrintError("Vault error: %v", err)
		return
	}

	if !v.Exists() {
		interactiveNoVault(v)
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

	interactiveMainLoop(v, masterPass, data)
}

// ── No Vault Flow ──

func interactiveNoVault(v *vault.Vault) {
	ui.PrintHeader("Welcome")
	ui.PrintWarn("No vault found. Create one to get started.")
	fmt.Println()
	printMenuOption("1", "➕ Add new credential")
	printMenuOption("2", "🔌 Quick connect (user@host)")
	printMenuOption("0", "🚪 Exit")
	fmt.Println()

	choice := ui.ReadLine(ui.InputPrompt("Select"))
	switch choice {
	case "1":
		addCommand.Run(addCommand, []string{})
	case "2":
		target := ui.ReadLine(ui.InputPrompt("Target (user@host)"))
		if target != "" {
			user, host, port := parseTarget(target)
			if user != "" && host != "" {
				ui.PrintConnecting(user, host, port)
				connectNewNoVault(v, user, host, port)
			} else {
				ui.PrintError("Invalid format. Use: user@host")
			}
		}
	}
}

// ── Main Interactive Loop ──

func interactiveMainLoop(v *vault.Vault, masterPass []byte, data *vault.VaultData) {
	for {
		fmt.Println()
		showConnectionList(data)
		showActionBar()

		choice := strings.TrimSpace(ui.ReadLine(
			fmt.Sprintf("\n  %s%s▸%s ", ui.Bold, ui.BrightGreen, ui.Reset)))
		if choice == "" {
			continue
		}

		switch strings.ToLower(choice) {
		case "q", "quit", "exit":
			fmt.Println()
			ui.PrintInfo("Goodbye!")
			fmt.Println()
			return

		case "a", "add":
			interactiveAdd(v, masterPass)
			data = reloadVault(v, masterPass, data)

		case "e", "edit":
			if len(data.Credentials) == 0 {
				ui.PrintWarn("No credentials to edit")
				continue
			}
			id := ui.ReadLine(ui.InputPrompt("# or name to edit"))
			alias := resolveAlias(data, id)
			if alias != "" {
				interactiveEdit(v, masterPass, data, alias)
				data = reloadVault(v, masterPass, data)
			}

		case "d", "del", "delete":
			if len(data.Credentials) == 0 {
				ui.PrintWarn("No credentials to delete")
				continue
			}
			id := ui.ReadLine(ui.InputPrompt("# or name to delete"))
			alias := resolveAlias(data, id)
			if alias != "" {
				interactiveDelete(v, masterPass, data, alias)
				data = reloadVault(v, masterPass, data)
			}

		case "m", "master":
			interactiveMaster(v, masterPass)

		default:
			// Try as number → direct connect
			if num, err := strconv.Atoi(choice); err == nil {
				if num >= 1 && num <= len(data.Credentials) {
					cred := &data.Credentials[num-1]
					connectFromMenu(cred)
				} else {
					ui.PrintError("Invalid number. Choose 1-%d", len(data.Credentials))
				}
			} else {
				// Search by name/host/user
				results := searchCredentials(data, choice)
				if len(results) == 0 {
					ui.PrintWarn("No match for '%s'", choice)
				} else if len(results) == 1 {
					cred := results[0]
					ui.PrintInfo("Found: %s%s%s (%s@%s)",
						ui.BrightWhite, cred.Alias, ui.Reset+ui.Cyan, cred.User, cred.Host)
					if ui.Confirm(fmt.Sprintf("  Connect to %s%s%s?",
						ui.BrightCyan, cred.Alias, ui.Reset)) {
						connectFromMenu(cred)
					}
				} else {
					showSearchResults(results, choice)
					pick := strings.TrimSpace(ui.ReadLine(
						fmt.Sprintf("  %s%s▸%s %s# to connect (Enter to cancel):%s ",
							ui.Bold, ui.BrightGreen, ui.Reset, ui.Dim, ui.Reset)))
					if pick != "" {
						if num, err := strconv.Atoi(pick); err == nil && num >= 1 && num <= len(results) {
							connectFromMenu(results[num-1])
						} else {
							ui.PrintError("Invalid selection")
						}
					}
				}
			}
		}
	}
}

// ── Connect from Menu ──

func connectFromMenu(cred *vault.Credential) {
	ui.PrintConnecting(cred.User, cred.Host, cred.Port)
	if err := connectWithCred(cred); err != nil {
		ui.PrintError("Connection failed: %v", err)
	}
	ui.FlushStdin()
	ui.PrintSessionEnd()
}

// ── Display Helpers ──

func showConnectionList(data *vault.VaultData) {
	if len(data.Credentials) == 0 {
		ui.PrintHeader("No Connections")
		ui.PrintWarn("No saved credentials yet. Press 'a' to add one.")
		return
	}

	ui.PrintHeader("Connections")

	headers := []string{"#", "NAME", "TARGET", "PORT", "AUTH"}
	rows := make([][]string, 0, len(data.Credentials))
	for i, c := range data.Credentials {
		authIcon := "🔑 pass"
		if c.AuthType == vault.AuthPrivateKey {
			authIcon = "📄 key"
		}
		rows = append(rows, []string{
			fmt.Sprintf("%d", i+1),
			c.Alias,
			fmt.Sprintf("%s@%s", c.User, c.Host),
			fmt.Sprintf("%d", c.Port),
			authIcon,
		})
	}

	ui.PrintTable(headers, rows)
}

func showActionBar() {
	fmt.Println()
	ui.PrintDivider()
	fmt.Printf("  %s[#]%s connect   %s[a]%sdd   %s[e]%sdit   %s[d]%sel   %s[q]%suit   %s· type to search%s\n",
		ui.BrightGreen, ui.Reset,
		ui.BrightCyan, ui.Reset,
		ui.BrightYellow, ui.Reset,
		ui.BrightRed, ui.Reset,
		ui.Dim, ui.Reset,
		ui.Dim+ui.BrightBlack, ui.Reset)
	ui.PrintDivider()
}

func showSearchResults(results []*vault.Credential, query string) {
	fmt.Println()
	ui.PrintHeader(fmt.Sprintf("Search: \"%s\"  ─  %d matches", query, len(results)))

	headers := []string{"#", "NAME", "TARGET", "PORT", "AUTH"}
	rows := make([][]string, 0, len(results))
	for i, c := range results {
		authIcon := "🔑 pass"
		if c.AuthType == vault.AuthPrivateKey {
			authIcon = "📄 key"
		}
		rows = append(rows, []string{
			fmt.Sprintf("%d", i+1),
			c.Alias,
			fmt.Sprintf("%s@%s", c.User, c.Host),
			fmt.Sprintf("%d", c.Port),
			authIcon,
		})
	}

	ui.PrintTable(headers, rows)
}

func printMenuOption(key, text string) {
	fmt.Printf("  %s%s[%s]%s %s%s%s\n",
		ui.Bold, ui.BrightCyan, key, ui.Reset,
		ui.BrightWhite, text, ui.Reset)
}

// ── Search ──

func searchCredentials(data *vault.VaultData, query string) []*vault.Credential {
	query = strings.ToLower(query)
	var results []*vault.Credential
	for i := range data.Credentials {
		c := &data.Credentials[i]
		target := strings.ToLower(c.User + "@" + c.Host)
		if strings.Contains(strings.ToLower(c.Alias), query) ||
			strings.Contains(strings.ToLower(c.Host), query) ||
			strings.Contains(strings.ToLower(c.User), query) ||
			strings.Contains(target, query) {
			results = append(results, c)
		}
	}
	return results
}

func resolveAlias(data *vault.VaultData, input string) string {
	input = strings.TrimSpace(input)
	if input == "" {
		return ""
	}
	// Try as number
	if num, err := strconv.Atoi(input); err == nil {
		if num >= 1 && num <= len(data.Credentials) {
			return data.Credentials[num-1].Alias
		}
		ui.PrintError("Invalid number. Choose 1-%d", len(data.Credentials))
		return ""
	}
	// Try as alias (case-insensitive)
	for _, c := range data.Credentials {
		if strings.EqualFold(c.Alias, input) {
			return c.Alias
		}
	}
	ui.PrintError("Credential '%s' not found", input)
	return ""
}

// ── Inline CRUD (no re-auth needed) ──

func interactiveAdd(v *vault.Vault, masterPass []byte) {
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

	defaultAlias := fmt.Sprintf("%s@%s", user, host)
	alias := ui.ReadLine(ui.InputPrompt(fmt.Sprintf("Name [%s]", defaultAlias)))
	if alias == "" {
		alias = defaultAlias
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
	ui.PrintSuccess("Credential '%s' saved!", alias)
}

func interactiveEdit(v *vault.Vault, masterPass []byte, data *vault.VaultData, alias string) {
	cred := v.FindByAlias(data, alias)
	if cred == nil {
		ui.PrintError("Credential '%s' not found", alias)
		return
	}

	ui.PrintHeader("Editing: " + alias)
	ui.PrintInfo("Press Enter to keep current value")
	fmt.Println()

	err := v.UpdateCredential(masterPass, alias, func(c *vault.Credential) {
		newHost := ui.ReadLine(ui.InputPrompt(fmt.Sprintf("Host [%s]", c.Host)))
		if newHost != "" {
			c.Host = newHost
		}

		newUser := ui.ReadLine(ui.InputPrompt(fmt.Sprintf("User [%s]", c.User)))
		if newUser != "" {
			c.User = newUser
		}

		newPort := ui.ReadLine(ui.InputPrompt(fmt.Sprintf("Port [%d]", c.Port)))
		if newPort != "" {
			if p, err := strconv.Atoi(newPort); err == nil {
				c.Port = p
			}
		}

		newAlias := ui.ReadLine(ui.InputPrompt(fmt.Sprintf("Name [%s]", c.Alias)))
		if newAlias != "" {
			c.Alias = newAlias
		}

		newAuthType := ui.ReadLine(ui.InputPrompt(fmt.Sprintf("Auth type [%s]", c.AuthType)))
		newAuthType = strings.ToLower(strings.TrimSpace(newAuthType))

		if newAuthType == "key" || newAuthType == "privatekey" {
			c.AuthType = vault.AuthPrivateKey
			keyPath := ui.ReadLine(ui.InputPrompt("Path to new private key (empty to keep)"))
			if keyPath != "" {
				keyData, err := os.ReadFile(keyPath)
				if err != nil {
					ui.PrintError("Error reading key: %v", err)
				} else {
					c.PrivateKey = string(keyData)
				}
			}
			if ui.Confirm("  Update key passphrase?") {
				keyPass, err := ui.ReadPassword(ui.PasswordPrompt("New key passphrase"))
				if err == nil {
					c.KeyPass = string(keyPass)
				}
			}
		} else if newAuthType == "password" || (newAuthType == "" && c.AuthType == vault.AuthPassword) {
			if ui.Confirm("  Update password?") {
				pass, err := ui.ReadPassword(ui.PasswordPrompt("New SSH password"))
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

	fmt.Println()
	ui.PrintSuccess("Credential updated!")
}

func interactiveDelete(v *vault.Vault, masterPass []byte, data *vault.VaultData, alias string) {
	cred := v.FindByAlias(data, alias)
	if cred == nil {
		ui.PrintError("Credential '%s' not found", alias)
		return
	}

	ui.PrintHeader("Delete Credential")
	ui.PrintKeyValue("Name:", cred.Alias)
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
}

func interactiveMaster(v *vault.Vault, masterPass []byte) {
	ui.PrintHeader("Change Master Password")

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
		ui.PrintError("Passwords don't match")
		return
	}

	if len(newPass) < 4 {
		ui.PrintError("Master password must be at least 4 characters")
		return
	}

	if err := v.ChangeMasterPassword(masterPass, newPass); err != nil {
		ui.PrintError("%v", err)
		return
	}

	fmt.Println()
	ui.PrintSuccess("Master password changed!")
}

func reloadVault(v *vault.Vault, masterPass []byte, fallback *vault.VaultData) *vault.VaultData {
	data, err := v.Load(masterPass)
	if err != nil {
		return fallback
	}
	return data
}
