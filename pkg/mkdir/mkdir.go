package mkdir

import (
	"fmt"
	"os"
)

// Options for mkdir
type Options struct {
	Parents bool // -p
}

// MakeDir creates directory path.
func MakeDir(path string, opts Options) error {
	if opts.Parents {
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("mkdir: cannot create directory '%s': %w", path, err)
		}
		return nil
	}
	if err := os.Mkdir(path, 0755); err != nil {
		return fmt.Errorf("mkdir: cannot create directory '%s': %w", path, err)
	}
	return nil
}
