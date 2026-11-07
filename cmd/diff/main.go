package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/user/win-tools/pkg/diff"
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
		fmt.Println("win-tools diff", version)
		os.Exit(0)
	}
	opts := diff.Options{}
	var files []string
	for _, arg := range expandShort(os.Args[1:]) {
		switch arg {
		case "-u", "--unified":
			opts.Unified = true
		case "-q", "--brief":
			opts.Brief = true
		case "-h", "--help":
			fmt.Fprintf(os.Stderr, "Usage: diff [-uq] FILE1 FILE2\n")
			os.Exit(0)
		default:
			files = append(files, arg)
		}
	}
	if len(files) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: diff [-uq] FILE1 FILE2\n")
		os.Exit(2)
	}
	differ, err := diff.Diff(files[0], files[1], opts, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	if differ {
		os.Exit(1)
	}
}