package main

import (
	"fmt"
	"os"

	"github.com/user/win-tools/pkg/which"
)

var version = "dev"

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Println("win-tools which", version)
		os.Exit(0)
	}

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: which [-a] COMMAND...\n")
		os.Exit(2)
	}

	showAll := false
	var commands []string

	for _, arg := range os.Args[1:] {
		switch arg {
		case "-a", "--all":
			showAll = true
		default:
			commands = append(commands, arg)
		}
	}

	for _, cmd := range commands {
		if showAll {
			results, err := which.WhichAll(cmd)
			if err != nil {
				fmt.Fprintf(os.Stderr, "which: %s: %v\n", cmd, err)
				os.Exit(1)
			}
			for _, r := range results {
				fmt.Println(r)
			}
		} else {
			result, err := which.Which(cmd)
			if err != nil {
				fmt.Fprintf(os.Stderr, "which: %s: %v\n", cmd, err)
				os.Exit(1)
			}
			fmt.Println(result)
		}
	}
}
