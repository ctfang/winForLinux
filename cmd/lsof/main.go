package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/user/win-tools/pkg/lsof"
)

var version = "dev"

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-V") {
		fmt.Println("win-tools lsof", version)
		os.Exit(0)
	}
	opts := lsof.Options{Internet: true}
	args := os.Args[1:]
	sawI := false
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "--help":
			fmt.Fprintf(os.Stderr, "Usage: lsof -i [TCP|UDP][:PORT] [-sTCP:LISTEN]\n")
			os.Exit(0)
		case a == "-i":
			sawI = true
			opts.Internet = true
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				i++
				proto, port, err := lsof.ParseIFilter(args[i])
				if err != nil {
					fmt.Fprintln(os.Stderr, "lsof:", err)
					os.Exit(2)
				}
				opts.Proto = proto
				opts.Port = port
			}
		case strings.HasPrefix(a, "-i") && len(a) > 2:
			sawI = true
			opts.Internet = true
			proto, port, err := lsof.ParseIFilter(a[2:])
			if err != nil {
				fmt.Fprintln(os.Stderr, "lsof:", err)
				os.Exit(2)
			}
			opts.Proto = proto
			opts.Port = port
		case strings.EqualFold(a, "-sTCP:LISTEN"):
			opts.Listen = true
		case strings.HasPrefix(a, "-"):
			fmt.Fprintf(os.Stderr, "lsof: invalid option: %s\n", a)
			os.Exit(2)
		default:
			proto, port, err := lsof.ParseIFilter(a)
			if err != nil {
				fmt.Fprintln(os.Stderr, "lsof:", err)
				os.Exit(2)
			}
			opts.Proto = proto
			opts.Port = port
		}
	}
	if !sawI && opts.Port == 0 && opts.Proto == "" && !opts.Listen {
		// allow bare lsof -i via default Internet; if no args at all, require -i
		if len(args) == 0 {
			fmt.Fprintf(os.Stderr, "Usage: lsof -i [TCP|UDP][:PORT] [-sTCP:LISTEN]\n")
			os.Exit(2)
		}
	}
	_ = sawI
	if err := lsof.Lsof(opts, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "lsof:", err)
		os.Exit(1)
	}
}