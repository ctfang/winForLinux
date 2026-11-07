package cp

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Options for cp
type Options struct {
	Recursive bool // -r / -R
	Force     bool // -f
}

// Copy copies src to dst.
func Copy(src, dst string, opts Options) error {
	info, err := os.Lstat(src)
	if err != nil {
		return fmt.Errorf("cp: cannot stat '%s': %w", src, err)
	}

	// If dst is existing directory, copy into it
	if di, err := os.Stat(dst); err == nil && di.IsDir() {
		dst = filepath.Join(dst, filepath.Base(src))
	}

	if info.IsDir() {
		if !opts.Recursive {
			return fmt.Errorf("cp: -r not specified; omitting directory '%s'", src)
		}
		return copyDir(src, dst, opts)
	}
	return copyFile(src, dst, opts)
}

func copyFile(src, dst string, opts Options) error {
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("cp: cannot open '%s': %w", src, err)
	}
	defer in.Close()

	if opts.Force {
		_ = os.Remove(dst)
	}

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("cp: cannot create '%s': %w", dst, err)
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return fmt.Errorf("cp: error copying '%s' to '%s': %w", src, dst, err)
	}
	return nil
}

func copyDir(src, dst string, opts Options) error {
	if err := os.MkdirAll(dst, 0755); err != nil {
		return fmt.Errorf("cp: cannot create directory '%s': %w", dst, err)
	}
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	for _, e := range entries {
		s := filepath.Join(src, e.Name())
		d := filepath.Join(dst, e.Name())
		if err := Copy(s, d, opts); err != nil {
			return err
		}
	}
	return nil
}
