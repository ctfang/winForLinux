package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/user/win-tools/pkg/sortcmd"
)

var version = "dev"

func expandShort(args []string) []string {
	var out []string
	for _, arg := range args {
		if strings.HasPrefix(arg, "--") {
			out = append(out, arg)
			continue
		}
		if strings.HasPrefix(arg, "-") && len(arg) > 2 && !strings.HasPrefix(arg, "-k") {
			all := true
			for _, ch := range arg[1:] {
				if (ch < 'a' || ch > 'z') && (ch < 'A' || ch > 'Z') {
					all = false
					break
				}
			}
			if all {
				for _, ch := range arg[1:] {
					out = append(out, "-"+string(ch))
				}
				continue
			}
		}
		out = append(out, arg)
	}
	return out
}

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-V") {
		fmt.Println("win-tools sort", version)
		os.Exit(0)
	}

	opts := sortcmd.Options{}
	var files []string
	args := expandShort(os.Args[1:])
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "-r" || arg == "--reverse":
			opts.Reverse = true
		case arg == "-n" || arg == "--numeric-sort":
			opts.Numeric = true
		case arg == "-u" || arg == "--unique":
			opts.Unique = true
		case arg == "-k" || arg == "--key":
			if i+1 >= len(args) {
				fmt.Fprintf(os.Stderr, "sort: option requires an argument -- k\n")
				os.Exit(2)
			}
			i++
			n, err := strconv.Atoi(args[i])
			if err != nil {
				fmt.Fprintf(os.Stderr, "sort: invalid key: %s\n", args[i])
				os.Exit(2)
			}
			opts.KeyField = n
		case strings.HasPrefix(arg, "-k"):
			n, err := strconv.Atoi(arg[2:])
			if err != nil {
				fmt.Fprintf(os.Stderr, "sort: invalid key: %s\n", arg[2:])
				os.Exit(2)
			}
			opts.KeyField = n
		case arg == "-h" || arg == "--help":
			fmt.Fprintf(os.Stderr, "Usage: sort [-rnu] [-k N] [FILE...]\n")
			os.Exit(0)
		default:
			files = append(files, arg)
		}
	}

	if len(files) == 0 {
		if err := sortcmd.SortLines(os.Stdin, opts, os.Stdout); err != nil {
			fmt.Fprintf(os.Stderr, "sort: %v\n", err)
			os.Exit(1)
		}
		return
	}
	for _, f := range files {
		file, err := os.Open(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "sort: %s: %v\n", f, err)
			continue
		}
		if err := sortcmd.SortLines(file, opts, os.Stdout); err != nil {
			fmt.Fprintf(os.Stderr, "sort: %s: %v\n", f, err)
		}
		file.Close()
	}
}
