package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/user/win-tools/pkg/sed"
)

var version = "dev"

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-V") {
		fmt.Println("win-tools sed", version)
		os.Exit(0)
	}
	opts := sed.Options{}
	var files []string
	args := os.Args[1:]
	exprFromArg := false
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "-h" || a == "--help":
			fmt.Fprintf(os.Stderr, "Usage: sed [-n] [-i] [-e script] [script] [FILE...]\n  basic: s/pat/repl/[gi]  d  p  q  addr/regex/\n")
			os.Exit(0)
		case a == "-n" || a == "--quiet" || a == "--silent":
			opts.Quiet = true
		case a == "-i" || a == "--in-place":
			opts.InPlace = true
		case a == "-e" || a == "--expression":
			i++
			if i >= len(args) {
				fmt.Fprintln(os.Stderr, "sed: option requires an argument -- e")
				os.Exit(2)
			}
			opts.Expressions = append(opts.Expressions, args[i])
			exprFromArg = true
		case strings.HasPrefix(a, "-e") && len(a) > 2:
			opts.Expressions = append(opts.Expressions, a[2:])
			exprFromArg = true
		case strings.HasPrefix(a, "-") && a != "-":
			fmt.Fprintf(os.Stderr, "sed: invalid option: %s\n", a)
			os.Exit(2)
		default:
			if !exprFromArg && len(opts.Expressions) == 0 {
				opts.Expressions = append(opts.Expressions, a)
				exprFromArg = true
			} else {
				files = append(files, a)
			}
		}
	}
	if len(opts.Expressions) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: sed [-n] [-i] [-e script] [script] [FILE...]")
		os.Exit(2)
	}
	if len(files) == 0 {
		if opts.InPlace {
			fmt.Fprintln(os.Stderr, "sed: -i requires a file")
			os.Exit(2)
		}
		if err := sed.Run(os.Stdin, opts, os.Stdout); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}
	exitCode := 0
	for _, f := range files {
		if err := sed.RunFile(f, opts); err != nil {
			fmt.Fprintln(os.Stderr, err)
			exitCode = 1
		}
	}
	os.Exit(exitCode)
}