package main

import (
	"fmt"
	"os"

	"github.com/user/win-tools/pkg/clear"
)

var version = "dev"

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--version", "-V":
			fmt.Println("win-tools clear", version)
			os.Exit(0)
		case "-h", "--help":
			fmt.Fprintf(os.Stderr, "Usage: clear\n")
			os.Exit(0)
		}
	}
	if err := clear.Clear(os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}