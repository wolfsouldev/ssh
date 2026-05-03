//go:build windows

package ui

// FlushStdin is a no-op on Windows.
func FlushStdin() {}
