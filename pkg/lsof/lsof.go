package lsof

import (
	"fmt"
	"io"
	"strings"

	"github.com/user/win-tools/pkg/winnet"
)

// Options for a minimal lsof -i subset.
type Options struct {
	Internet bool // -i
	Port     int
	Proto    string // tcp/udp/empty
	Listen   bool
}

// Lsof prints internet sockets similar to lsof -i.
func Lsof(opts Options, w io.Writer) error {
	nopts := winnet.Options{
		TCP:       opts.Proto == "" || opts.Proto == "tcp" || opts.Proto == "TCP",
		UDP:       opts.Proto == "" || opts.Proto == "udp" || opts.Proto == "UDP",
		Listening: opts.Listen,
		All:       !opts.Listen,
		WithProc:  true,
		LocalPort: opts.Port,
	}
	conns, err := winnet.List(nopts)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "%-16s %8s %6s %s\n", "COMMAND", "PID", "TYPE", "NAME")
	for _, c := range conns {
		cmd := c.ProcName
		if cmd == "" {
			cmd = "-"
		}
		name := winnet.FormatAddr(c.LocalIP, c.LocalPort)
		if c.Proto == "tcp" || c.Proto == "tcp6" {
			if c.StateCode == 2 {
				name += " (LISTEN)"
			} else {
				name += "->" + winnet.FormatAddr(c.RemoteIP, c.RemotePort) + " (" + c.State + ")"
			}
		}
		typ := strings.ToUpper(c.Proto)
		fmt.Fprintf(w, "%-16s %8d %6s %s\n", trunc(cmd, 16), c.PID, typ, name)
	}
	return nil
}

func trunc(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}

// ParseIFilter parses -i arguments like TCP:8080, :8080, 8080, TCP.
func ParseIFilter(s string) (proto string, port int, err error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return "", 0, nil
	}
	upper := strings.ToUpper(s)
	if upper == "TCP" || upper == "UDP" {
		return strings.ToLower(upper), 0, nil
	}
	if strings.HasPrefix(s, ":") {
		_, err = fmt.Sscanf(s, ":%d", &port)
		return "", port, err
	}
	if strings.Contains(s, ":") {
		parts := strings.SplitN(s, ":", 2)
		proto = strings.ToLower(parts[0])
		if parts[1] != "" {
			_, err = fmt.Sscanf(parts[1], "%d", &port)
		}
		return proto, port, err
	}
	_, err = fmt.Sscanf(s, "%d", &port)
	return "", port, err
}