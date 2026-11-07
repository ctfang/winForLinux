package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/user/win-tools/pkg/cat"
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
			allFlags := true
			for _, ch := range arg[1:] {
				if (ch < 'a' || ch > 'z') && (ch < 'A' || ch > 'Z') {
					allFlags = false
					break
				}
			}
			if allFlags {
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
		fmt.Println("win-tools cat", version)
		os.Exit(0)
	}

	opts := cat.Options{}
	var files []string

	for _, arg := range expandShort(os.Args[1:]) {
		switch arg {
		case "-n", "--number":
			opts.Number = true
		case "-E", "--show-ends":
			opts.ShowEnds = true
		case "-T", "--show-tabs":
			opts.ShowTabs = true
		case "-s", "--squeeze-blank":
			opts.SqueezeBlank = true
		case "-h", "--help":
			fmt.Fprintf(os.Stderr, "Usage: cat [OPTIONS] [FILE...]\n  -n number lines\n  -E show ends\n  -T show tabs\n  -s squeeze blank\n")
			os.Exit(0)
		default:
			files = append(files, arg)
		}
	}

	if len(files) > 0 {
		for _, filename := range files {
			func() {
				file, err := os.Open(filename)
				if err != nil {
					fmt.Fprintf(os.Stderr, "cat: %s: %v\n", filename, err)
					return
				}
				defer file.Close()
				if err := cat.Cat(file, opts, os.Stdout); err != nil {
					fmt.Fprintf(os.Stderr, "cat: %s: %v\n", filename, err)
				}
			}()
		}
	} else {
		if err := cat.Cat(os.Stdin, opts, os.Stdout); err != nil {
			fmt.Fprintf(os.Stderr, "cat: %v\n", err)
		}
	}
}
