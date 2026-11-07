package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/user/win-tools/pkg/touch"
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
		fmt.Println("win-tools touch", version)
		os.Exit(0)
	}

	opts := touch.Options{}
	var files []string
	for _, arg := range expandShort(os.Args[1:]) {
		switch arg {
		case "-a":
			opts.AccessOnly = true
		case "-m":
			opts.ModOnly = true
		case "-c", "--no-create":
			opts.NoCreate = true
		case "-h", "--help":
			fmt.Fprintf(os.Stderr, "Usage: touch [-amc] FILE...\n")
			os.Exit(0)
		default:
			files = append(files, arg)
		}
	}
	if len(files) == 0 {
		fmt.Fprintf(os.Stderr, "Usage: touch [-amc] FILE...\n")
		os.Exit(2)
	}
	exitCode := 0
	for _, f := range files {
		if err := touch.Touch(f, opts); err != nil {
			fmt.Fprintln(os.Stderr, err)
			exitCode = 1
		}
	}
	os.Exit(exitCode)
}
