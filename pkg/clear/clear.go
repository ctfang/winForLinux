package clear

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
)

// Clear clears the terminal screen.
func Clear(w io.Writer) error {
	// ANSI clear screen + move cursor home (works in modern Windows consoles)
	if _, err := io.WriteString(w, "\033[H\033[2J\033[3J"); err == nil {
		if f, ok := w.(*os.File); ok {
			_ = f.Sync()
		}
	}
	if runtime.GOOS == "windows" {
		// Also invoke cls for hosts that ignore ANSI
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = w
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			// ANSI already attempted; ignore cls failure
			_ = err
		}
	}
	return nil
}

// ClearANSI only emits ANSI sequences (for tests / redirected output).
func ClearANSI(w io.Writer) error {
	_, err := fmt.Fprint(w, "\033[H\033[2J\033[3J")
	return err
}