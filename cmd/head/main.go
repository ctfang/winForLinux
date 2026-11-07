package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/user/win-tools/pkg/head"
)

var version = "dev"

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-V") {
		fmt.Println("win-tools head", version)
		os.Exit(0)
	}

	opts := head.Options{Lines: 10}
	var files []string
	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "-h" || arg == "--help":
			fmt.Fprintf(os.Stderr, "Usage: head [-n N|-c N] [FILE...]\n")
			os.Exit(0)
		case arg == "-n" || arg == "--lines":
			if i+1 >= len(args) {
				fmt.Fprintf(os.Stderr, "head: option requires an argument -- n\n")
				os.Exit(2)
			}
			i++
			n, err := strconv.Atoi(args[i])
			if err != nil {
				fmt.Fprintf(os.Stderr, "head: invalid number: %s\n", args[i])
				os.Exit(2)
			}
			opts.Lines = n
			opts.Bytes = 0
		case strings.HasPrefix(arg, "-n"):
			n, err := strconv.Atoi(arg[2:])
			if err != nil {
				fmt.Fprintf(os.Stderr, "head: invalid number: %s\n", arg[2:])
				os.Exit(2)
			}
			opts.Lines = n
			opts.Bytes = 0
		case arg == "-c" || arg == "--bytes":
			if i+1 >= len(args) {
				fmt.Fprintf(os.Stderr, "head: option requires an argument -- c\n")
				os.Exit(2)
			}
			i++
			n, err := strconv.Atoi(args[i])
			if err != nil {
				fmt.Fprintf(os.Stderr, "head: invalid number: %s\n", args[i])
				os.Exit(2)
			}
			opts.Bytes = n
			opts.Lines = 0
		case len(arg) > 1 && arg[0] == '-' && arg[1] >= '0' && arg[1] <= '9':
			n, err := strconv.Atoi(arg[1:])
			if err != nil {
				fmt.Fprintf(os.Stderr, "head: invalid number: %s\n", arg[1:])
				os.Exit(2)
			}
			opts.Lines = n
			opts.Bytes = 0
		default:
			files = append(files, arg)
		}
	}

	run := func(name string, r *os.File) {
		if err := head.Head(r, opts, os.Stdout); err != nil {
			fmt.Fprintf(os.Stderr, "head: %s: %v\n", name, err)
		}
	}

	if len(files) == 0 {
		run("stdin", os.Stdin)
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
			fmt.Fprintf(os.Stderr, "head: %s: %v\n", f, err)
			continue
		}
		run(f, file)
		file.Close()
	}
}
