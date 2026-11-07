package touch

import (
	"fmt"
	"os"
	"time"
)

// Options for touch
type Options struct {
	AccessOnly bool // -a
	ModOnly    bool // -m
	NoCreate   bool // -c
}

// Touch creates file or updates timestamps.
func Touch(path string, opts Options) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		if opts.NoCreate {
			return nil
		}
		f, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("touch: cannot create '%s': %w", path, err)
		}
		return f.Close()
	}
	if err != nil {
		return fmt.Errorf("touch: %s: %w", path, err)
	}

	now := time.Now()
	atime, mtime := now, now
	info, _ := os.Stat(path)
	if info != nil {
		if opts.AccessOnly && !opts.ModOnly {
			mtime = info.ModTime()
		}
		if opts.ModOnly && !opts.AccessOnly {
			// keep atime approximately as now is fine on Windows; Chtimes sets both
			atime = now
			mtime = now
		}
	}
	if err := os.Chtimes(path, atime, mtime); err != nil {
		return fmt.Errorf("touch: %s: %w", path, err)
	}
	return nil
}
