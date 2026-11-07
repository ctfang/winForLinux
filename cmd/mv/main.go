package main

import (
	"fmt"
	"os"

	"github.com/user/win-tools/pkg/mv"
)

var version = "dev"

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-V") {
		fmt.Println("win-tools mv", version)
		os.Exit(0)
	}
	args := []string{}
	for _, a := range os.Args[1:] {
		if a == "-h" || a == "--help" {
			fmt.Fprintf(os.Stderr, "Usage: mv SOURCE... DEST\n")
			os.Exit(0)
		}
		if a == "-f" || a == "--force" {
			continue
		}
		args = append(args, a)
	}
	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: mv SOURCE... DEST\n")
		os.Exit(2)
	}
	dst := args[len(args)-1]
	srcs := args[:len(args)-1]
	if len(srcs) > 1 {
		info, err := os.Stat(dst)
		if err != nil || !info.IsDir() {
			fmt.Fprintf(os.Stderr, "mv: target '%s' is not a directory\n", dst)
			os.Exit(1)
		}
	}
	exitCode := 0
	for _, s := range srcs {
		if err := mv.Move(s, dst); err != nil {
			fmt.Fprintln(os.Stderr, err)
			exitCode = 1
		}
	}
	os.Exit(exitCode)
}
