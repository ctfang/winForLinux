package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/user/win-tools/pkg/find"
)

var version = "dev"

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-V") {
		fmt.Println("win-tools find", version)
		os.Exit(0)
	}
	opts := find.Options{MaxDepth: -1}
	var paths []string
	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "-h" || a == "--help":
			fmt.Fprintf(os.Stderr, "Usage: find [PATH...] [EXPRESSION]\n  -name PATTERN  -iname PATTERN  -type f|d|l\n  -maxdepth N  -mindepth N  -size N[cKMGb]  -mtime N\n")
			os.Exit(0)
		case a == "-name":
			i++
			if i >= len(args) {
				die("find: missing argument to -name")
			}
			opts.Name = args[i]
		case a == "-iname":
			i++
			if i >= len(args) {
				die("find: missing argument to -iname")
			}
			opts.IName = args[i]
		case a == "-type":
			i++
			if i >= len(args) {
				die("find: missing argument to -type")
			}
			opts.Type = args[i]
		case a == "-maxdepth":
			i++
			if i >= len(args) {
				die("find: missing argument to -maxdepth")
			}
			n, err := strconv.Atoi(args[i])
			if err != nil {
				die("find: invalid -maxdepth")
			}
			opts.MaxDepth = n
		case a == "-mindepth":
			i++
			if i >= len(args) {
				die("find: missing argument to -mindepth")
			}
			n, err := strconv.Atoi(args[i])
			if err != nil {
				die("find: invalid -mindepth")
			}
			opts.MinDepth = n
		case a == "-size":
			i++
			if i >= len(args) {
				die("find: missing argument to -size")
			}
			opts.Size = args[i]
		case a == "-mtime":
			i++
			if i >= len(args) {
				die("find: missing argument to -mtime")
			}
			opts.MTime = args[i]
		case a == "-print":
			// default
		case len(a) > 0 && a[0] == '-':
			die("find: unknown predicate: " + a)
		default:
			paths = append(paths, a)
		}
	}
	if err := find.Find(paths, opts, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func die(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(2)
}