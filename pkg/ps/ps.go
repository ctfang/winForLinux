package ps

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/user/win-tools/pkg/winproc"
)

// Options for ps
type Options struct {
	All      bool  // -e / -A
	Full     bool  // -f
	PIDs     []uint32 // -p
}

// Ps lists processes.
func Ps(opts Options, w io.Writer) error {
	procs, err := winproc.List()
	if err != nil {
		return err
	}

	if len(opts.PIDs) > 0 {
		want := map[uint32]bool{}
		for _, p := range opts.PIDs {
			want[p] = true
		}
		filtered := procs[:0]
		for _, p := range procs {
			if want[p.PID] {
				filtered = append(filtered, p)
			}
		}
		procs = filtered
	}

	sort.Slice(procs, func(i, j int) bool { return procs[i].PID < procs[j].PID })

	if opts.Full {
		fmt.Fprintf(w, "%-8s %8s %8s %6s %s\n", "UID", "PID", "PPID", "THR", "CMD")
		for _, p := range procs {
			cmd := p.Name
			if p.ExePath != "" {
				cmd = p.ExePath
			}
			fmt.Fprintf(w, "%-8s %8d %8d %6d %s\n", "-", p.PID, p.PPID, p.Threads, cmd)
		}
		return nil
	}

	fmt.Fprintf(w, "%8s %s\n", "PID", "CMD")
	for _, p := range procs {
		fmt.Fprintf(w, "%8d %s\n", p.PID, p.Name)
	}
	return nil
}

// ParsePIDList parses comma-separated PIDs.
func ParsePIDList(s string) ([]uint32, error) {
	var out []uint32
	for _, part := range strings.Split(s, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		var id uint32
		_, err := fmt.Sscanf(part, "%d", &id)
		if err != nil {
			return nil, fmt.Errorf("invalid pid: %s", part)
		}
		out = append(out, id)
	}
	return out, nil
}