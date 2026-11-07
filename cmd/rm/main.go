package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/user/win-tools/pkg/rm"
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
		fmt.Println("win-tools rm", version)
		os.Exit(0)
	}
	opts := rm.Options{}
	var paths []string
	for _, arg := range expandShort(os.Args[1:]) {
		switch arg {
		case "-r", "-R", "--recursive":
			opts.Recursive = true
		case "-f", "--force":
			opts.Force = true
		case "-h", "--help":
			fmt.Fprintf(os.Stderr, "Usage: rm [-rf] FILE...\n")
			os.Exit(0)
		default:
			paths = append(paths, arg)
		}
	}
	if len(paths) == 0 {
		fmt.Fprintf(os.Stderr, "Usage: rm [-rf] FILE...\n")
		os.Exit(2)
	}
	exitCode := 0
	for _, p := range paths {
		if err := rm.Remove(p, opts); err != nil {
			fmt.Fprintln(os.Stderr, err)
			exitCode = 1
		}
	}
	os.Exit(exitCode)
}
