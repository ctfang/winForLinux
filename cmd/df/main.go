package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/user/win-tools/pkg/df"
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
		fmt.Println("win-tools df", version)
		os.Exit(0)
	}
	opts := df.Options{}
	var paths []string
	for _, arg := range expandShort(os.Args[1:]) {
		switch arg {
		case "-h", "--human-readable":
			opts.Human = true
		case "--help":
			fmt.Fprintf(os.Stderr, "Usage: df [-h] [FILE...]\n")
			os.Exit(0)
		default:
			if strings.HasPrefix(arg, "-") {
				fmt.Fprintf(os.Stderr, "df: invalid option: %s\n", arg)
				os.Exit(2)
			}
			paths = append(paths, arg)
		}
	}
	if err := df.Df(paths, opts, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}