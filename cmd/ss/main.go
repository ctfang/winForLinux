package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/user/win-tools/pkg/ss"
)

var version = "dev"

func expandShort(args []string) []string {
	var out []string
	for _, arg := range args {
		if strings.HasPrefix(arg, "--") {
			out = append(out, arg)
			continue
		}
		if strings.HasPrefix(arg, "-") && len(arg) > 2 {
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
		fmt.Println("win-tools ss", version)
		os.Exit(0)
	}
	opts := ss.Options{All: true, Numeric: true}
	args := expandShort(os.Args[1:])
	for _, a := range args {
		switch {
		case a == "-t", a == "--tcp":
			opts.TCP = true
		case a == "-u", a == "--udp":
			opts.UDP = true
		case a == "-l", a == "--listening":
			opts.Listening = true
			opts.All = false
		case a == "-a", a == "--all":
			opts.All = true
			opts.Listening = false
		case a == "-n", a == "--numeric":
			opts.Numeric = true
		case a == "-p", a == "--processes":
			opts.Processes = true
		case a == "--help":
			fmt.Fprintf(os.Stderr, "Usage: ss [-tulnap] [PORT|:PORT]\n")
			os.Exit(0)
		case strings.HasPrefix(a, "-"):
			fmt.Fprintf(os.Stderr, "ss: invalid option: %s\n", a)
			os.Exit(2)
		default:
			port, err := parsePortArg(a)
			if err != nil {
				fmt.Fprintln(os.Stderr, "ss:", err)
				os.Exit(2)
			}
			opts.LocalPort = port
		}
	}
	if !opts.TCP && !opts.UDP {
		opts.TCP, opts.UDP = true, true
	}
	if err := ss.Ss(opts, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "ss:", err)
		os.Exit(1)
	}
}

func parsePortArg(s string) (int, error) {
	// sport = :22  or  :22  or  22
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, " ", "")
	lower := strings.ToLower(s)
	if strings.HasPrefix(lower, "sport=") {
		s = s[len("sport="):]
	}
	s = strings.TrimPrefix(s, ":")
	n, err := strconv.Atoi(s)
	if err != nil || n <= 0 || n > 65535 {
		return 0, fmt.Errorf("invalid port filter: %s", s)
	}
	return n, nil
}