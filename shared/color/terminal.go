package color

import (
	"os"
	"syscall"
	"unsafe"
)

var (
	kernel32           = syscall.NewLazyDLL("kernel32.dll")
	procGetConsoleMode = kernel32.NewProc("GetConsoleMode")
)

// IsTerminal checks whether stdout is connected to a terminal.
func IsTerminal() bool {
	return isTerminal(os.Stdout.Fd())
}

func isTerminal(fd uintptr) bool {
	var st uint32
	r, _, _ := procGetConsoleMode.Call(fd, uintptr(unsafe.Pointer(&st)))
	return r != 0
}

// SupportsColor checks whether the terminal supports ANSI colors.
func SupportsColor() bool {
	// Windows Terminal and ConEmu support ANSI
	if os.Getenv("WT_SESSION") != "" {
		return true // Windows Terminal
	}
	if os.Getenv("ConEmuANSI") == "ON" {
		return true // ConEmu
	}
	// Traditional cmd needs ENABLE_VIRTUAL_TERMINAL_PROCESSING
	return IsTerminal()
}
