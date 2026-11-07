package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/user/win-tools/pkg/ps"
)

var version = "dev"

func expandShort(args []string) []string {
	var out []string
	for _, arg := range args {
		if strings.HasPrefix(arg, "--") {
			out = append(out, arg)
			continue
		}
		if strings.HasPrefix(arg, "-") && len(arg) > 2 && !strings.HasPrefix(arg, "-p") {
			all := true
			for _, ch := range arg[1:] {
				if (ch < 'a' || ch > 'z') && (ch < 'A' || ch > 'Z') {
					all = false
					break
				}
			}
			if all {
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
		fmt.Println("win-tools ps", version)
		os.Exit(0)
	}
	opts := ps.Options{All: true}
	args := expandShort(os.Args[1:])
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "-e", a == "-A", a == "--all":
			opts.All = true
		case a == "-f", a == "--full":
			opts.Full = true
		case a == "-ef":
			opts.All, opts.Full = true, true
		case a == "-p", a == "--pid":
			i++
			if i >= len(args) {
				fmt.Fprintln(os.Stderr, "ps: option requires an argument -- p")
				os.Exit(2)
			}
			pids, err := ps.ParsePIDList(args[i])
			if err != nil {
				fmt.Fprintln(os.Stderr, "ps:", err)
				os.Exit(2)
			}
			opts.PIDs = append(opts.PIDs, pids...)
		case strings.HasPrefix(a, "-p") && len(a) > 2:
			pids, err := ps.ParsePIDList(a[2:])
			if err != nil {
				fmt.Fprintln(os.Stderr, "ps:", err)
				os.Exit(2)
			}
			opts.PIDs = append(opts.PIDs, pids...)
		case a == "--help":
			fmt.Fprintf(os.Stderr, "Usage: ps [-e|-A] [-f] [-p PID[,PID...]]\n")
			os.Exit(0)
		case strings.HasPrefix(a, "-"):
			// combined -ef already handled; unknown single after expand
			fmt.Fprintf(os.Stderr, "ps: invalid option: %s\n", a)
			os.Exit(2)
		default:
			fmt.Fprintf(os.Stderr, "ps: unexpected argument: %s\n", a)
			os.Exit(2)
		}
	}
	if err := ps.Ps(opts, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "ps:", err)
		os.Exit(1)
	}
}