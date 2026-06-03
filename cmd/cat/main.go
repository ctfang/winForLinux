package main

import (
	"fmt"
	"os"

	"github.com/user/win-tools/pkg/cat"
)

var version = "dev"

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Println("win-tools cat", version)
		os.Exit(0)
	}

	opts := cat.Options{}
	var files []string

	// Parse options
	for _, arg := range os.Args[1:] {
		switch arg {
		case "-n", "--number":
			opts.Number = true
		case "-E", "--show-ends":
			opts.ShowEnds = true
		case "-T", "--show-tabs":
			opts.ShowTabs = true
		case "-s", "--squeeze-blank":
			opts.SqueezeBlank = true
		default:
			files = append(files, arg)
		}
	}

	// Read from files or stdin
	if len(files) > 0 {
		for _, filename := range files {
			file, err := os.Open(filename)
			if err != nil {
				fmt.Fprintf(os.Stderr, "cat: %s: %v\n", filename, err)
				continue
			}
			defer file.Close()
			cat.Cat(file, opts, os.Stdout)
		}
	} else {
		cat.Cat(os.Stdin, opts, os.Stdout)
	}
}
