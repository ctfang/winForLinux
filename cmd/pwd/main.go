package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/user/win-tools/pkg/pwd"
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
		fmt.Println("win-tools pwd", version)
		os.Exit(0)
	}
	opts := pwd.Options{}
	for _, arg := range expandShort(os.Args[1:]) {
		switch arg {
		case "-L", "--logical":
			opts.Physical = false
		case "-P", "--physical":
			opts.Physical = true
		case "-h", "--help":
			fmt.Fprintf(os.Stderr, "Usage: pwd [-LP]\n")
			os.Exit(0)
		default:
			fmt.Fprintf(os.Stderr, "pwd: invalid option: %s\n", arg)
			os.Exit(2)
		}
	}
	if err := pwd.Pwd(opts, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}