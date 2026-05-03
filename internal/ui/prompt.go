package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"
)

// ── ANSI Color Codes ──

const (
	Reset     = "\033[0m"
	Bold      = "\033[1m"
	Dim       = "\033[2m"
	Italic    = "\033[3m"
	Underline = "\033[4m"
	Blink     = "\033[5m"

	// Foreground
	Black   = "\033[30m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"

	// Bright Foreground
	BrightBlack   = "\033[90m"
	BrightRed     = "\033[91m"
	BrightGreen   = "\033[92m"
	BrightYellow  = "\033[93m"
	BrightBlue    = "\033[94m"
	BrightMagenta = "\033[95m"
	BrightCyan    = "\033[96m"
	BrightWhite   = "\033[97m"

	// Background
	BgBlack  = "\033[40m"
	BgRed    = "\033[41m"
	BgGreen  = "\033[42m"
	BgYellow = "\033[43m"
	BgBlue   = "\033[44m"
	BgCyan   = "\033[46m"
)

// ── ASCII Art Banner ──

const Banner = `
` + BrightCyan + `   ███████╗███████╗██╗  ██╗██╗  ██╗` + Reset + `
` + Cyan + `   ██╔════╝██╔════╝██║  ██║██║  ██║` + Reset + `
` + BrightGreen + `   ███████╗███████╗███████║███████║` + Reset + `
` + Green + `   ╚════██║╚════██║██╔══██║██╔══██║` + Reset + `
` + BrightCyan + `   ███████║███████║██║  ██║██║  ██║` + Reset + `
` + Cyan + `   ╚══════╝╚══════╝╚═╝  ╚═╝╚═╝  ╚═╝` + Reset + `
` + Dim + BrightGreen + `   ─── Secure Shell Handler ───────` + Reset + `
`

// ── Astronaut ASCII Art ──

const Astronaut = `` +
	Dim + `          ·  .  ` + BrightYellow + `★` + Reset + Dim + `  .  ·` + Reset + "\n" +
	BrightWhite + `         ╭───────────╮` + Reset + "\n" +
	BrightWhite + `         │ ╭───────╮ │` + Reset + "\n" +
	BrightWhite + `         │ │` + BrightGreen + ` ◉   ◉ ` + BrightWhite + `│ │` + Reset + "\n" +
	BrightWhite + `         │ │` + Dim + `   ─   ` + Reset + BrightWhite + `│ │` + Reset + "\n" +
	BrightWhite + `         │ ╰───────╯ │` + Reset + "\n" +
	BrightWhite + `         ╰─────┬─────╯` + Reset + "\n" +
	Cyan + `         ╭─────┴─────╮` + Reset + "\n" +
	Cyan + `         │` + BrightGreen + ` >` + Bold + ` SSHH ` + Reset + Cyan + `_ │` + Reset + "\n" +
	Cyan + `         ╰───────────╯` + Reset + "\n" +
	Dim + `          ·  .  ` + BrightBlue + `·` + Reset + Dim + `  .  ·` + Reset + "\n"

// PrintBanner prints the ASCII art banner.
func PrintBanner() {
	fmt.Print(Banner)
}

// PrintAstronaut prints the astronaut ASCII art.
func PrintAstronaut() {
	fmt.Print(Astronaut)
}

// ── Input Functions ──

// ReadLine reads a line of input from the user.
func ReadLine(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}

// ReadPassword reads a password from the user without echoing.
func ReadPassword(prompt string) ([]byte, error) {
	fmt.Print(prompt)
	pass, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return nil, fmt.Errorf("reading password: %w", err)
	}
	return pass, nil
}

// Confirm asks the user a yes/no question.
func Confirm(prompt string) bool {
	answer := ReadLine(prompt + Yellow + " [y/N]: " + Reset)
	return strings.ToLower(answer) == "y" || strings.ToLower(answer) == "yes"
}

// ── Output Functions ──

// PrintSuccess prints a success message in green.
func PrintSuccess(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("  %s%s[✓]%s %s%s%s\n", Bold, BrightGreen, Reset, Green, msg, Reset)
}

// PrintError prints an error message in red.
func PrintError(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("  %s%s[✗]%s %s%s%s\n", Bold, BrightRed, Reset, Red, msg, Reset)
}

// PrintInfo prints an info message in cyan.
func PrintInfo(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("  %s%s[>]%s %s%s%s\n", Bold, BrightCyan, Reset, Cyan, msg, Reset)
}

// PrintWarn prints a warning message in yellow.
func PrintWarn(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("  %s%s[!]%s %s%s%s\n", Bold, BrightYellow, Reset, Yellow, msg, Reset)
}

// PrintHeader prints a styled header with box drawing.
func PrintHeader(text string) {
	width := len(text) + 4
	top := "╔" + strings.Repeat("═", width) + "╗"
	bot := "╚" + strings.Repeat("═", width) + "╝"
	fmt.Printf("\n  %s%s%s\n", BrightCyan, top, Reset)
	fmt.Printf("  %s║%s  %s%s%s  %s║%s\n", BrightCyan, Reset, Bold+BrightWhite, text, Reset, BrightCyan, Reset)
	fmt.Printf("  %s%s%s\n\n", BrightCyan, bot, Reset)
}

// PrintDivider prints a horizontal divider line.
func PrintDivider() {
	fmt.Printf("  %s%s%s\n", Dim+BrightBlack, strings.Repeat("─", 50), Reset)
}

// PrintKeyValue prints a label-value pair with colors.
func PrintKeyValue(label, value string) {
	fmt.Printf("  %s%s%-12s%s %s%s%s\n", Bold, BrightCyan, label, Reset, BrightWhite, value, Reset)
}

// PrintConnecting prints the connecting animation header.
func PrintConnecting(user, host string, port int) {
	fmt.Println()
	PrintDivider()
	fmt.Printf("  %s%s⚡ CONNECTING%s\n", Bold, BrightGreen, Reset)
	PrintDivider()
	PrintKeyValue("Target:", fmt.Sprintf("%s@%s", user, host))
	PrintKeyValue("Port:", fmt.Sprintf("%d", port))
	PrintDivider()
	fmt.Println()
}

// PrintSessionEnd prints a message when the SSH session ends.
func PrintSessionEnd() {
	fmt.Println()
	PrintDivider()
	fmt.Printf("  %s%s⏏  SESSION CLOSED%s\n", Bold, BrightYellow, Reset)
	PrintDivider()
	fmt.Println()
}

// ── Table Rendering ──

// PrintTable prints a hacker-styled table.
func PrintTable(headers []string, rows [][]string) {
	if len(rows) == 0 {
		PrintInfo("No entries found.")
		return
	}

	// Calculate column widths
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = len(h)
	}
	for _, row := range rows {
		for i, cell := range row {
			if i < len(widths) && len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	// Calculate total table width
	totalWidth := 1 // left border
	for _, w := range widths {
		totalWidth += w + 3 // cell + padding + separator
	}

	// Top border
	fmt.Printf("  %s┌", Dim+BrightCyan)
	for i, w := range widths {
		fmt.Print(strings.Repeat("─", w+2))
		if i < len(widths)-1 {
			fmt.Print("┬")
		}
	}
	fmt.Printf("┐%s\n", Reset)

	// Header row
	fmt.Printf("  %s│%s", Dim+BrightCyan, Reset)
	for i, h := range headers {
		fmt.Printf(" %s%s%-*s%s ", Bold, BrightGreen, widths[i], h, Reset)
		fmt.Printf("%s│%s", Dim+BrightCyan, Reset)
	}
	fmt.Println()

	// Header separator
	fmt.Printf("  %s├", Dim+BrightCyan)
	for i, w := range widths {
		fmt.Print(strings.Repeat("─", w+2))
		if i < len(widths)-1 {
			fmt.Print("┼")
		}
	}
	fmt.Printf("┤%s\n", Reset)

	// Data rows
	for _, row := range rows {
		fmt.Printf("  %s│%s", Dim+BrightCyan, Reset)
		for i, cell := range row {
			if i < len(widths) {
				fmt.Printf(" %s%-*s%s ", BrightWhite, widths[i], cell, Reset)
				fmt.Printf("%s│%s", Dim+BrightCyan, Reset)
			}
		}
		fmt.Println()
	}

	// Bottom border
	fmt.Printf("  %s└", Dim+BrightCyan)
	for i, w := range widths {
		fmt.Print(strings.Repeat("─", w+2))
		if i < len(widths)-1 {
			fmt.Print("┴")
		}
	}
	fmt.Printf("┘%s\n", Reset)
}

// ── Prompt Styling ──

// PasswordPrompt returns a styled password prompt string.
func PasswordPrompt(label string) string {
	return fmt.Sprintf("  %s%s🔐 %s:%s ", Bold, BrightYellow, label, Reset)
}

// InputPrompt returns a styled input prompt string.
func InputPrompt(label string) string {
	return fmt.Sprintf("  %s%s▸%s %s%s:%s ", Bold, BrightCyan, Reset, Dim, label, Reset)
}
