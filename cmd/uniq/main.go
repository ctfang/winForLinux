package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/user/win-tools/pkg/uniq"
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
		fmt.Println("win-tools uniq", version)
		os.Exit(0)
	}
	opts := uniq.Options{}
	var files []string
	for _, arg := range expandShort(os.Args[1:]) {
		switch arg {
		case "-c", "--count":
			opts.Count = true
		case "-d", "--repeated":
			opts.Repeated = true
		case "-u", "--unique":
			opts.UniqueOnly = true
		case "-h", "--help":
			fmt.Fprintf(os.Stderr, "Usage: uniq [-cdu] [FILE]\n")
			os.Exit(0)
		default:
			files = append(files, arg)
		}
	}
	if len(files) == 0 {
		if err := uniq.Uniq(os.Stdin, opts, os.Stdout); err != nil {
			fmt.Fprintf(os.Stderr, "uniq: %v\n", err)
			os.Exit(1)
		}
		return
	}
	f, err := os.Open(files[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "uniq: %s: %v\n", files[0], err)
		os.Exit(1)
	}
	defer f.Close()
	if err := uniq.Uniq(f, opts, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "uniq: %v\n", err)
		os.Exit(1)
	}
}
