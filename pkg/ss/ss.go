package ss

import (
	"fmt"
	"io"

	"github.com/user/win-tools/pkg/winnet"
)

// Options mirrors common Linux ss flags.
type Options struct {
	TCP       bool
	UDP       bool
	Listening bool
	All       bool
	Numeric   bool
	Processes bool
	LocalPort int
}

// Ss prints socket statistics in ss-like format.
func Ss(opts Options, w io.Writer) error {
	nopts := winnet.Options{
		TCP:       opts.TCP,
		UDP:       opts.UDP,
		Listening: opts.Listening,
		All:       opts.All || !opts.Listening,
		WithProc:  opts.Processes,
		LocalPort: opts.LocalPort,
	}
	if !opts.TCP && !opts.UDP {
		nopts.TCP, nopts.UDP = true, true
	}
	conns, err := winnet.List(nopts)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "Netid  State      Local Address:Port     Peer Address:Port      Process\n")
	for _, c := range conns {
		local := winnet.FormatAddr(c.LocalIP, c.LocalPort)
		peer := "*:*"
		if (c.Proto == "tcp" || c.Proto == "tcp6") && c.StateCode != 2 {
			peer = winnet.FormatAddr(c.RemoteIP, c.RemotePort)
		}
		state := c.State
		if c.Proto == "udp" || c.Proto == "udp6" {
			state = "UNCONN"
		}
		proc := ""
		if opts.Processes && c.PID != 0 {
			if c.ProcName != "" {
				proc = fmt.Sprintf("users:((\"%s\",pid=%d))", c.ProcName, c.PID)
			} else {
				proc = fmt.Sprintf("pid=%d", c.PID)
			}
		}
		netid := c.Proto
		fmt.Fprintf(w, "%-6s %-10s %-22s %-22s %s\n", netid, state, local, peer, proc)
	}
	return nil
}