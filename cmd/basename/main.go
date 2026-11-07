package main

import (
	"fmt"
	"os"

	"github.com/user/win-tools/pkg/basename"
)

var version = "dev"

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-V") {
		fmt.Println("win-tools basename", version)
		os.Exit(0)
	}
	args := os.Args[1:]
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		fmt.Fprintf(os.Stderr, "Usage: basename NAME [SUFFIX]\n       basename -a NAME...\n")
		if len(args) == 0 {
			os.Exit(2)
		}
		os.Exit(0)
	}
	if args[0] == "-a" || args[0] == "--multiple" {
		for _, a := range args[1:] {
			fmt.Println(basename.Base(a, ""))
		}
		return
	}
	suffix := ""
	if len(args) >= 2 {
		suffix = args[1]
	}
	fmt.Println(basename.Base(args[0], suffix))
}
