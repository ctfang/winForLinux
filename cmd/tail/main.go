package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/user/win-tools/pkg/tail"
)

var version = "dev"

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-V") {
		fmt.Println("win-tools tail", version)
		os.Exit(0)
	}

	opts := tail.Options{Lines: 10}
	var files []string
	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "-h" || arg == "--help":
			fmt.Fprintf(os.Stderr, "Usage: tail [-n N] [-f] [FILE...]\n")
			os.Exit(0)
		case arg == "-f" || arg == "--follow":
			opts.Follow = true
		case arg == "-n" || arg == "--lines":
			if i+1 >= len(args) {
				fmt.Fprintf(os.Stderr, "tail: option requires an argument -- n\n")
				os.Exit(2)
			}
			i++
			n, err := strconv.Atoi(args[i])
			if err != nil {
				fmt.Fprintf(os.Stderr, "tail: invalid number: %s\n", args[i])
				os.Exit(2)
			}
			opts.Lines = n
		case strings.HasPrefix(arg, "-n"):
			n, err := strconv.Atoi(arg[2:])
			if err != nil {
				fmt.Fprintf(os.Stderr, "tail: invalid number: %s\n", arg[2:])
				os.Exit(2)
			}
			opts.Lines = n
		case len(arg) > 1 && arg[0] == '-' && arg[1] >= '0' && arg[1] <= '9':
			n, err := strconv.Atoi(arg[1:])
			if err != nil {
				fmt.Fprintf(os.Stderr, "tail: invalid number: %s\n", arg[1:])
				os.Exit(2)
			}
			opts.Lines = n
		default:
			files = append(files, arg)
		}
	}

	if len(files) == 0 {
		if err := tail.Tail(os.Stdin, opts, os.Stdout); err != nil {
			fmt.Fprintf(os.Stderr, "tail: %v\n", err)
			os.Exit(1)
		}
		return
	}

	for i, f := range files {
		if len(files) > 1 {
			if i > 0 {
				fmt.Println()
			}
			fmt.Printf("==> %s <==\n", f)
		}
		file, err := os.Open(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "tail: %s: %v\n", f, err)
			continue
		}
		if err := tail.Tail(file, opts, os.Stdout); err != nil {
			fmt.Fprintf(os.Stderr, "tail: %s: %v\n", f, err)
		}
		file.Close()
	}

	if opts.Follow && len(files) == 1 {
		if err := tail.FollowFile(files[0], os.Stdout); err != nil {
			fmt.Fprintf(os.Stderr, "tail: %v\n", err)
			os.Exit(1)
		}
	}
}
