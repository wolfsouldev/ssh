package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"
)

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
	answer := ReadLine(prompt + " [y/N]: ")
	return strings.ToLower(answer) == "y" || strings.ToLower(answer) == "yes"
}

// PrintSuccess prints a success message in green.
func PrintSuccess(format string, args ...interface{}) {
	fmt.Printf("  ✓ "+format+"\n", args...)
}

// PrintError prints an error message in red.
func PrintError(format string, args ...interface{}) {
	fmt.Printf("  ✗ "+format+"\n", args...)
}

// PrintInfo prints an info message.
func PrintInfo(format string, args ...interface{}) {
	fmt.Printf("  → "+format+"\n", args...)
}

// PrintHeader prints a header.
func PrintHeader(text string) {
	fmt.Printf("\n  ═══ %s ═══\n\n", text)
}

// PrintTable prints a simple table.
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

	// Print header
	fmt.Print("  ")
	for i, h := range headers {
		fmt.Printf("%-*s  ", widths[i], h)
	}
	fmt.Println()

	// Print separator
	fmt.Print("  ")
	for i := range headers {
		fmt.Print(strings.Repeat("─", widths[i]))
		fmt.Print("  ")
	}
	fmt.Println()

	// Print rows
	for _, row := range rows {
		fmt.Print("  ")
		for i, cell := range row {
			if i < len(widths) {
				fmt.Printf("%-*s  ", widths[i], cell)
			}
		}
		fmt.Println()
	}
}
