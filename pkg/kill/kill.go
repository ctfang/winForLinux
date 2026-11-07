package kill

import (
	"fmt"
	"os/exec"
	"strconv"

	"github.com/user/win-tools/pkg/winproc"
)

// Options for kill
type Options struct {
	Signal int // -9 maps to forceful terminate
}

// Kill terminates processes by PID.
func Kill(pids []uint32, opts Options) error {
	force := opts.Signal == 9
	var msgs []string
	for _, pid := range pids {
		if err := winproc.Kill(pid, force); err != nil {
			// Fallback: Windows taskkill often succeeds where TerminateProcess is denied (UWP etc.)
			if ferr := taskKill(pid, force || opts.Signal == 9); ferr != nil {
				msgs = append(msgs, fmt.Sprintf("(%d) - %v", pid, err))
			}
		}
	}
	if len(msgs) > 0 {
		return fmt.Errorf("kill: %s", join(msgs, "; "))
	}
	return nil
}

func taskKill(pid uint32, force bool) error {
	args := []string{"/PID", strconv.FormatUint(uint64(pid), 10)}
	if force {
		args = append(args, "/F")
	}
	cmd := exec.Command("taskkill", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("taskkill: %v (%s)", err, trim(string(out)))
	}
	return nil
}

func trim(s string) string {
	for len(s) > 0 && (s[len(s)-1] == '\n' || s[len(s)-1] == '\r' || s[len(s)-1] == ' ') {
		s = s[:len(s)-1]
	}
	return s
}

func join(a []string, sep string) string {
	if len(a) == 0 {
		return ""
	}
	out := a[0]
	for i := 1; i < len(a); i++ {
		out += sep + a[i]
	}
	return out
}