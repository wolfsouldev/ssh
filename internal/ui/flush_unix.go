//go:build !windows

package ui

import (
	"os"
	"syscall"
)

// FlushStdin discards any pending input from stdin.
// This is critical after an SSH interactive session where the terminal
// was in raw mode — leftover bytes can cause subsequent reads to hang.
func FlushStdin() {
	fd := int(os.Stdin.Fd())
	if err := syscall.SetNonblock(fd, true); err != nil {
		return
	}
	buf := make([]byte, 512)
	for {
		n, _ := syscall.Read(fd, buf)
		if n <= 0 {
			break
		}
	}
	_ = syscall.SetNonblock(fd, false)
}
