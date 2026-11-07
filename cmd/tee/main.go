package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/user/win-tools/pkg/tee"
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
		fmt.Println("win-tools tee", version)
		os.Exit(0)
	}
	opts := tee.Options{}
	var files []string
	for _, arg := range expandShort(os.Args[1:]) {
		switch arg {
		case "-a", "--append":
			opts.Append = true
		case "-h", "--help":
			fmt.Fprintf(os.Stderr, "Usage: tee [-a] [FILE...]\n")
			os.Exit(0)
		default:
			files = append(files, arg)
		}
	}
	if err := tee.Tee(os.Stdin, files, opts, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}