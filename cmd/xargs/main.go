package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/user/win-tools/pkg/xargs"
)

var version = "dev"

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-V") {
		fmt.Println("win-tools xargs", version)
		os.Exit(0)
	}
	opts := xargs.Options{}
	args := os.Args[1:]
	var command []string
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "-h" || a == "--help":
			fmt.Fprintf(os.Stderr, "Usage: xargs [-0] [-n N] [-I repl] [-d delim] [-r] [command [args...]]\n")
			os.Exit(0)
		case a == "-0" || a == "--null":
			opts.NullSep = true
		case a == "-r" || a == "--no-run-if-empty":
			opts.NoRunIfEmpty = true
		case a == "-n" || a == "--max-args":
			i++
			if i >= len(args) {
				fmt.Fprintln(os.Stderr, "xargs: option requires an argument -- n")
				os.Exit(2)
			}
			n, err := strconv.Atoi(args[i])
			if err != nil || n < 1 {
				fmt.Fprintln(os.Stderr, "xargs: invalid number")
				os.Exit(2)
			}
			opts.MaxArgs = n
		case strings.HasPrefix(a, "-n") && len(a) > 2:
			n, err := strconv.Atoi(a[2:])
			if err != nil || n < 1 {
				fmt.Fprintln(os.Stderr, "xargs: invalid number")
				os.Exit(2)
			}
			opts.MaxArgs = n
		case a == "-I" || a == "--replace":
			i++
			if i >= len(args) {
				fmt.Fprintln(os.Stderr, "xargs: option requires an argument -- I")
				os.Exit(2)
			}
			opts.Replace = args[i]
		case a == "-i":
			opts.Replace = "{}"
		case a == "-d" || a == "--delimiter":
			i++
			if i >= len(args) {
				fmt.Fprintln(os.Stderr, "xargs: option requires an argument -- d")
				os.Exit(2)
			}
			opts.Delimiter = args[i]
		case a == "--":
			command = args[i+1:]
			i = len(args)
		case strings.HasPrefix(a, "-"):
			fmt.Fprintf(os.Stderr, "xargs: invalid option: %s\n", a)
			os.Exit(2)
		default:
			command = args[i:]
			i = len(args)
		}
	}
	if err := xargs.Run(os.Stdin, command, opts, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}