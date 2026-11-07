package pwd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Options for pwd
type Options struct {
	Physical bool // -P resolve symlinks
}

// Pwd writes the working directory.
func Pwd(opts Options, w io.Writer) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("pwd: %w", err)
	}
	if opts.Physical {
		if abs, err := filepath.Abs(dir); err == nil {
			if resolved, err := filepath.EvalSymlinks(abs); err == nil {
				dir = resolved
			} else {
				dir = abs
			}
		}
	}
	_, err = fmt.Fprintln(w, dir)
	return err
}