package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/user/win-tools/pkg/killall"
)

var version = "dev"

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-V") {
		fmt.Println("win-tools killall", version)
		os.Exit(0)
	}
	opts := killall.Options{Signal: 15}
	var names []string
	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "--help":
			fmt.Fprintf(os.Stderr, "Usage: killall [-9|-s SIGNAL] [-q] [-v] NAME...\n")
			os.Exit(0)
		case a == "-9", a == "-KILL", a == "-kill":
			opts.Signal = 9
		case a == "-q", a == "--quiet":
			opts.Quiet = true
		case a == "-v", a == "--verbose":
			opts.Verbose = true
		case a == "-s", a == "--signal":
			i++
			if i >= len(args) {
				fmt.Fprintln(os.Stderr, "killall: option requires an argument -- s")
				os.Exit(2)
			}
			sig, err := parseSignal(args[i])
			if err != nil {
				fmt.Fprintln(os.Stderr, "killall:", err)
				os.Exit(2)
			}
			opts.Signal = sig
		case strings.HasPrefix(a, "-") && len(a) > 1:
			if n, err := strconv.Atoi(a[1:]); err == nil {
				opts.Signal = n
				continue
			}
			fmt.Fprintf(os.Stderr, "killall: invalid option: %s\n", a)
			os.Exit(2)
		default:
			names = append(names, a)
		}
	}
	if len(names) == 0 {
		fmt.Fprintf(os.Stderr, "Usage: killall [-9|-s SIGNAL] [-q] [-v] NAME...\n")
		os.Exit(2)
	}
	if err := killall.Killall(names, opts); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func parseSignal(s string) (int, error) {
	s = strings.TrimPrefix(strings.ToUpper(s), "SIG")
	switch s {
	case "KILL", "9":
		return 9, nil
	case "TERM", "15":
		return 15, nil
	default:
		n, err := strconv.Atoi(s)
		if err != nil {
			return 0, fmt.Errorf("unknown signal: %s", s)
		}
		return n, nil
	}
}