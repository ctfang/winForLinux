package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/user/win-tools/pkg/wfl"
)

var version = "dev"

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--version", "-V":
			fmt.Println("win-tools wfl", version)
			os.Exit(0)
		case "-h", "--help":
			printHelp()
			os.Exit(0)
		case "-c":
			if len(os.Args) < 3 {
				fmt.Fprintf(os.Stderr, "wfl: -c requires a command string\n")
				os.Exit(2)
			}
			line := strings.Join(os.Args[2:], " ")
			os.Exit(wfl.RunLine(line, os.Stdin, os.Stdout, os.Stderr))
		}
	}

	if len(os.Args) < 2 {
		printHelp()
		os.Exit(2)
	}

	// Join remaining args as a command line (operators passed as separate args work too)
	line := strings.Join(os.Args[1:], " ")
	os.Exit(wfl.RunLine(line, os.Stdin, os.Stdout, os.Stderr))
}

func printHelp() {
	fmt.Fprintf(os.Stderr, `wfl - winForLinux command runner (Linux-style combining)

Usage:
  wfl -c "COMMAND"
  wfl COMMAND [ARGS...]

Operators:
  |    pipe stdout of left to stdin of right
  &&   run right only if left succeeded
  ||   run right only if left failed
  ;    run sequentially
  &    run left in background

Examples:
  wfl -c "echo hello | grep hello"
  wfl -c "false || echo fallback"
  wfl -c "true && echo ok"
  wfl -c "echo bg & echo fg"

Note: In PowerShell, native pipes already work between win-tools .exe files:
  cat README.md | grep win
Use wfl when you need Linux-style && || & ; semantics that PowerShell treats differently.
`)
}
