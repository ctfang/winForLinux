package main

import (
	"fmt"
	"os"

	"github.com/user/win-tools/pkg/dirname"
)

var version = "dev"

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-V") {
		fmt.Println("win-tools dirname", version)
		os.Exit(0)
	}
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: dirname NAME\n")
		os.Exit(2)
	}
	if os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Fprintf(os.Stderr, "Usage: dirname NAME\n")
		os.Exit(0)
	}
	fmt.Println(dirname.Dir(os.Args[1]))
}
