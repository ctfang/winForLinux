package netstat

import (
	"fmt"
	"io"

	"github.com/user/win-tools/pkg/winnet"
)

// Options mirrors common Linux netstat flags.
type Options struct {
	TCP       bool
	UDP       bool
	Listening bool
	All       bool
	Numeric   bool
	Programs  bool
	LocalPort int
}

// Netstat prints connection table.
func Netstat(opts Options, w io.Writer) error {
	nopts := winnet.Options{
		TCP:       opts.TCP,
		UDP:       opts.UDP,
		Listening: opts.Listening,
		All:       opts.All || (!opts.Listening),
		WithProc:  opts.Programs,
		LocalPort: opts.LocalPort,
	}
	if !opts.TCP && !opts.UDP {
		nopts.TCP, nopts.UDP = true, true
	}
	conns, err := winnet.List(nopts)
	if err != nil {
		return err
	}

	if opts.Programs {
		fmt.Fprintf(w, "%-6s %-23s %-23s %-12s %s\n", "Proto", "Local Address", "Foreign Address", "State", "PID/Program")
	} else {
		fmt.Fprintf(w, "%-6s %-23s %-23s %s\n", "Proto", "Local Address", "Foreign Address", "State")
	}
	for _, c := range conns {
		local := winnet.FormatAddr(c.LocalIP, c.LocalPort)
		remote := "*"
		if c.Proto == "tcp" || c.Proto == "tcp6" {
			if c.StateCode == 2 { // LISTEN
				remote = "*:*"
			} else {
				remote = winnet.FormatAddr(c.RemoteIP, c.RemotePort)
			}
		} else {
			remote = "*:*"
		}
		state := c.State
		if c.Proto == "udp" || c.Proto == "udp6" {
			state = ""
		}
		if opts.Programs {
			prog := ""
			if c.PID != 0 {
				if c.ProcName != "" {
					prog = fmt.Sprintf("%d/%s", c.PID, c.ProcName)
				} else {
					prog = fmt.Sprintf("%d", c.PID)
				}
			}
			fmt.Fprintf(w, "%-6s %-23s %-23s %-12s %s\n", c.Proto, local, remote, state, prog)
		} else {
			fmt.Fprintf(w, "%-6s %-23s %-23s %s\n", c.Proto, local, remote, state)
		}
	}
	return nil
}