package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/user/win-tools/pkg/du"
)

var version = "dev"

func expandShort(args []string) []string {
	var out []string
	for _, arg := range args {
		if strings.HasPrefix(arg, "--") {
			out = append(out, arg)
			continue
		}
		if strings.HasPrefix(arg, "-") && len(arg) > 2 {
			// -dN special
			if arg[1] == 'd' && len(arg) > 2 && arg[2] >= '0' && arg[2] <= '9' {
				out = append(out, "-d", arg[2:])
				continue
			}
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
		fmt.Println("win-tools du", version)
		os.Exit(0)
	}
	opts := du.Options{}
	var paths []string
	args := expandShort(os.Args[1:])
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "-h" || a == "--human-readable":
			opts.Human = true
		case a == "-s" || a == "--summarize":
			opts.Summarize = true
		case a == "-a" || a == "--all":
			opts.All = true
		case a == "-d" || a == "--max-depth":
			i++
			if i >= len(args) {
				fmt.Fprintln(os.Stderr, "du: option requires an argument -- d")
				os.Exit(2)
			}
			n, err := strconv.Atoi(args[i])
			if err != nil {
				fmt.Fprintln(os.Stderr, "du: invalid max depth")
				os.Exit(2)
			}
			opts.MaxDepth = n
			opts.DepthSet = true
		case strings.HasPrefix(a, "--max-depth="):
			n, err := strconv.Atoi(a[len("--max-depth="):])
			if err != nil {
				fmt.Fprintln(os.Stderr, "du: invalid max depth")
				os.Exit(2)
			}
			opts.MaxDepth = n
			opts.DepthSet = true
		case a == "--help":
			fmt.Fprintf(os.Stderr, "Usage: du [-h] [-s] [-a] [-d N] [FILE...]\n")
			os.Exit(0)
		case strings.HasPrefix(a, "-"):
			fmt.Fprintf(os.Stderr, "du: invalid option: %s\n", a)
			os.Exit(2)
		default:
			paths = append(paths, a)
		}
	}
	if err := du.Du(paths, opts, os.Stdout); err != nil {
		os.Exit(1)
	}
}