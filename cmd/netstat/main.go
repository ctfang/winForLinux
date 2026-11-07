package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/user/win-tools/pkg/netstat"
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
		fmt.Println("win-tools netstat", version)
		os.Exit(0)
	}
	opts := netstat.Options{All: true, Numeric: true}
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
		case a == "-p", a == "--programs":
			opts.Programs = true
		case a == "-o": // Windows-style show owning PID
			opts.Programs = true
		case a == "--help":
			fmt.Fprintf(os.Stderr, "Usage: netstat [-tulnap] [PORT|:PORT]\n")
			os.Exit(0)
		case strings.HasPrefix(a, "-"):
			fmt.Fprintf(os.Stderr, "netstat: invalid option: %s\n", a)
			os.Exit(2)
		default:
			port, err := parsePort(a)
			if err != nil {
				fmt.Fprintln(os.Stderr, "netstat:", err)
				os.Exit(2)
			}
			opts.LocalPort = port
		}
	}
	if !opts.TCP && !opts.UDP {
		opts.TCP, opts.UDP = true, true
	}
	if err := netstat.Netstat(opts, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "netstat:", err)
		os.Exit(1)
	}
}

func parsePort(s string) (int, error) {
	s = strings.TrimPrefix(s, ":")
	n, err := strconv.Atoi(s)
	if err != nil || n <= 0 || n > 65535 {
		return 0, fmt.Errorf("invalid port: %s", s)
	}
	return n, nil
}