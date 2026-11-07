package rm

import (
	"fmt"
	"os"
)

// Options for rm
type Options struct {
	Recursive bool // -r
	Force     bool // -f
}

// Remove removes path.
func Remove(path string, opts Options) error {
	info, err := os.Lstat(path)
	if err != nil {
		if opts.Force && os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("rm: cannot remove '%s': %w", path, err)
	}
	if info.IsDir() {
		if !opts.Recursive {
			return fmt.Errorf("rm: cannot remove '%s': Is a directory", path)
		}
		if err := os.RemoveAll(path); err != nil {
			if opts.Force {
				return nil
			}
			return fmt.Errorf("rm: cannot remove '%s': %w", path, err)
		}
		return nil
	}
	if err := os.Remove(path); err != nil {
		if opts.Force {
			return nil
		}
		return fmt.Errorf("rm: cannot remove '%s': %w", path, err)
	}
	return nil
}
