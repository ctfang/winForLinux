package color

import "fmt"

// ANSI color codes
const (
	Reset   = 0
	Bold    = 1
	Red     = 31
	Green   = 32
	Yellow  = 33
	Blue    = 34
	Magenta = 35
	Cyan    = 36
	White   = 37
)

// Colorize wraps text with ANSI color codes
func Colorize(text string, code int) string {
	return fmt.Sprintf("\033[%dm%s\033[0m", code, text)
}

// BoldText makes text bold
func BoldText(text string) string {
	return fmt.Sprintf("\033[1m%s\033[0m", text)
}
