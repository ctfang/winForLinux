package main

import (
	"os"

	"github.com/user/win-tools/pkg/ls"
)

var version = "dev"

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		println("win-tools ls", version)
		os.Exit(0)
	}

	if err := ls.Run(os.Args[1:]); err != nil {
		os.Exit(1)
	}
}
