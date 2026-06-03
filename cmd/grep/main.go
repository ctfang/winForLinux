package main

import (
	"fmt"
	"os"

	"github.com/user/win-tools/pkg/grep"
	"github.com/user/win-tools/shared/color"
)

var version = "dev"

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Println("win-tools grep", version)
		os.Exit(0)
	}

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: grep [OPTIONS] PATTERN [FILE...]\n")
		os.Exit(2)
	}

	opts := grep.Options{}
	var pattern string
	var files []string

	// Parse options and extract pattern
	for _, arg := range os.Args[1:] {
		switch arg {
		case "-i", "--ignore-case":
			opts.IgnoreCase = true
		case "-n", "--line-number":
			opts.LineNumber = true
		case "-c", "--count":
			opts.Count = true
		case "-v", "--invert-match":
			opts.Invert = true
		case "-o", "--only-matching":
			opts.OnlyMatch = true
		case "-w", "--word-regexp":
			opts.WordRegexp = true
		case "-E", "--extended-regexp":
			// regexp.Compile already handles extended syntax; flag is accepted but unused
		case "-F", "--fixed-strings":
			opts.FixedString = true
		case "--color":
			opts.Color = resolveColorMode("always")
		case "--color=auto":
			opts.Color = resolveColorMode("auto")
		case "--color=always":
			opts.Color = resolveColorMode("always")
		case "--color=never":
			opts.Color = "never"
		default:
			if pattern == "" {
				pattern = arg
			} else {
				files = append(files, arg)
			}
		}
	}

	if pattern == "" {
		fmt.Fprintf(os.Stderr, "Usage: grep [OPTIONS] PATTERN [FILE...]\n")
		os.Exit(2)
	}

	// Read from files or stdin
	var exitCode int
	if len(files) > 0 {
		for _, filename := range files {
			file, err := os.Open(filename)
			if err != nil {
				fmt.Fprintf(os.Stderr, "grep: %s: %v\n", filename, err)
				exitCode = 2
				continue
			}
			defer file.Close()
			if err := grep.Grep(file, pattern, opts, os.Stdout); err != nil {
				fmt.Fprintf(os.Stderr, "grep: %s: %v\n", filename, err)
				exitCode = 2
			}
		}
	} else {
		if err := grep.Grep(os.Stdin, pattern, opts, os.Stdout); err != nil {
			fmt.Fprintf(os.Stderr, "grep: %v\n", err)
			exitCode = 2
		}
	}

	os.Exit(exitCode)
}

// resolveColorMode maps a color setting to an effective mode.
// "auto" resolves to "always" when stdout is a terminal that supports color,
// and "" (disabled) otherwise.
func resolveColorMode(mode string) string {
	if mode == "auto" {
		if color.SupportsColor() {
			return "always"
		}
		return ""
	}
	return mode
}
