package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/user/win-tools/pkg/echo"
)

var version = "dev"

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-V") {
		fmt.Println("win-tools echo", version)
		os.Exit(0)
	}
	opts := echo.Options{}
	raw := os.Args[1:]
	i := 0
	for ; i < len(raw); i++ {
		a := raw[i]
		if a == "--" {
			i++
			break
		}
		if !strings.HasPrefix(a, "-") || a == "-" {
			break
		}
		switch a {
		case "-n":
			opts.NoNewline = true
		case "-e":
			opts.EnableEsc = true
		case "-E":
			opts.EnableEsc = false
		case "-h", "--help":
			fmt.Fprintf(os.Stderr, "Usage: echo [-neE] [STRING...]\n")
			os.Exit(0)
		default:
			goto rest
		}
	}
rest:
	args := raw[i:]
	if err := echo.Echo(args, opts, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}