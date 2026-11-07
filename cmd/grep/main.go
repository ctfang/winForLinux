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

	args := grep.ExpandShortFlags(os.Args[1:])
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Usage: grep [OPTIONS] PATTERN [FILE...]\n")
		os.Exit(2)
	}

	opts := grep.Options{}
	var pattern string
	var files []string
	patternSet := false

	for i := 0; i < len(args); i++ {
		arg := args[i]
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
			// accepted; Go regexp is already extended-like
		case "-F", "--fixed-strings":
			opts.FixedString = true
		case "-r", "-R", "--recursive":
			opts.Recursive = true
		case "--color":
			opts.Color = resolveColorMode("always")
		case "--color=auto":
			opts.Color = resolveColorMode("auto")
		case "--color=always":
			opts.Color = resolveColorMode("always")
		case "--color=never":
			opts.Color = "never"
		case "-e", "--regexp":
			if i+1 >= len(args) {
				fmt.Fprintf(os.Stderr, "grep: option requires an argument -- e\n")
				os.Exit(2)
			}
			i++
			pattern = args[i]
			patternSet = true
		case "-h", "--help":
			fmt.Fprintf(os.Stderr, "Usage: grep [OPTIONS] PATTERN [FILE...]\n")
			fmt.Fprintf(os.Stderr, "  -i  ignore case\n  -n  line numbers\n  -c  count\n  -v  invert\n  -r  recursive\n  -w  word\n  -F  fixed string\n")
			os.Exit(0)
		default:
			if !patternSet {
				pattern = arg
				patternSet = true
			} else {
				files = append(files, arg)
			}
		}
	}

	if !patternSet || pattern == "" {
		fmt.Fprintf(os.Stderr, "Usage: grep [OPTIONS] PATTERN [FILE...]\n")
		os.Exit(2)
	}

	exitCode := 1 // 1 = no match (GNU grep convention)
	foundAny := false

	if opts.Recursive {
		roots := files
		if len(roots) == 0 {
			roots = []string{"."}
		}
		for _, root := range roots {
			found, err := grep.GrepRecursive(root, pattern, opts, os.Stdout)
			if err != nil {
				fmt.Fprintf(os.Stderr, "grep: %v\n", err)
				exitCode = 2
				continue
			}
			if found {
				foundAny = true
			}
		}
	} else if len(files) > 0 {
		multi := len(files) > 1
		for _, filename := range files {
			info, err := os.Stat(filename)
			if err != nil {
				fmt.Fprintf(os.Stderr, "grep: %s: %v\n", filename, err)
				exitCode = 2
				continue
			}
			if info.IsDir() {
				fmt.Fprintf(os.Stderr, "grep: %s: Is a directory\n", filename)
				exitCode = 2
				continue
			}
			file, err := os.Open(filename)
			if err != nil {
				fmt.Fprintf(os.Stderr, "grep: %s: %v\n", filename, err)
				exitCode = 2
				continue
			}
			name := ""
			if multi {
				name = filename
			}
			found, err := grep.GrepNamed(file, name, pattern, opts, os.Stdout)
			file.Close()
			if err != nil {
				fmt.Fprintf(os.Stderr, "grep: %s: %v\n", filename, err)
				exitCode = 2
				continue
			}
			if found {
				foundAny = true
			}
		}
	} else {
		found, err := grep.Grep(os.Stdin, pattern, opts, os.Stdout)
		if err != nil {
			fmt.Fprintf(os.Stderr, "grep: %v\n", err)
			os.Exit(2)
		}
		foundAny = found
	}

	if foundAny && exitCode != 2 {
		exitCode = 0
	}
	os.Exit(exitCode)
}

func resolveColorMode(mode string) string {
	if mode == "auto" {
		if color.SupportsColor() {
			return "always"
		}
		return ""
	}
	return mode
}
