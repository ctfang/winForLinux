package main

import (
	"fmt"
	"os"

	"github.com/user/win-tools/pkg/grep"
)

var version = "dev"

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Println("win-tools grep", version)
		os.Exit(0)
	}

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: grep [OPTIONS] PATTERN [FILE...]\n")
		os.Exit(2)
	}

	opts := grep.Options{
		Pattern: os.Args[1],
	}

	// Parse options
	var files []string
	for _, arg := range os.Args[2:] {
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
			opts.ExtendedReg = true
		case "-F", "--fixed-strings":
			opts.FixedString = true
		case "--color=auto":
			opts.Color = "auto"
		case "--color=always":
			opts.Color = "always"
		case "--color=never":
			opts.Color = "never"
		default:
			files = append(files, arg)
		}
	}

	// Read from files or stdin
	if len(files) > 0 {
		for _, filename := range files {
			file, err := os.Open(filename)
			if err != nil {
				fmt.Fprintf(os.Stderr, "grep: %s: %v\n", filename, err)
				continue
			}
			grep.Grep(file, opts.Pattern, opts, os.Stdout)
			file.Close()
		}
	} else {
		grep.Grep(os.Stdin, opts.Pattern, opts, os.Stdout)
	}
}
