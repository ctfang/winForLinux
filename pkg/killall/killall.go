package killall

import (
	"fmt"
	"strings"

	"github.com/user/win-tools/pkg/kill"
	"github.com/user/win-tools/pkg/winproc"
)

// Options for killall
type Options struct {
	Signal  int
	Quiet   bool
	Verbose bool
}

// Killall terminates processes by name.
func Killall(names []string, opts Options) error {
	var errs []string
	for _, name := range names {
		procs, err := winproc.ByName(name)
		if err != nil {
			errs = append(errs, err.Error())
			continue
		}
		if len(procs) == 0 {
			if !opts.Quiet {
				errs = append(errs, fmt.Sprintf("%s: no process found", name))
			}
			continue
		}
		for _, p := range procs {
			if err := kill.Kill([]uint32{p.PID}, kill.Options{Signal: opts.Signal}); err != nil {
				errs = append(errs, fmt.Sprintf("%s(%d): %v", p.Name, p.PID, err))
				continue
			}
			if opts.Verbose {
				fmt.Printf("Killed %s(%d)\n", p.Name, p.PID)
			}
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("killall: %s", strings.Join(errs, "; "))
	}
	return nil
}