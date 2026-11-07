package mv

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Move renames/moves src to dst.
func Move(src, dst string) error {
	if di, err := os.Stat(dst); err == nil && di.IsDir() {
		dst = filepath.Join(dst, filepath.Base(src))
	}
	if err := os.Rename(src, dst); err == nil {
		return nil
	} else {
		// Cross-device fallback
		info, serr := os.Lstat(src)
		if serr != nil {
			return fmt.Errorf("mv: cannot move '%s' to '%s': %w", src, dst, err)
		}
		if info.IsDir() {
			if err := copyDir(src, dst); err != nil {
				return fmt.Errorf("mv: cannot move '%s' to '%s': %w", src, dst, err)
			}
			if err := os.RemoveAll(src); err != nil {
				return fmt.Errorf("mv: moved but failed to remove '%s': %w", src, err)
			}
			return nil
		}
		if err := copyFile(src, dst, info.Mode()); err != nil {
			return fmt.Errorf("mv: cannot move '%s' to '%s': %w", src, dst, err)
		}
		if err := os.Remove(src); err != nil {
			return fmt.Errorf("mv: moved but failed to remove '%s': %w", src, err)
		}
		return nil
	}
}

func copyFile(src, dst string, mode os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}

func copyDir(src, dst string) error {
	if err := os.MkdirAll(dst, 0755); err != nil {
		return err
	}
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	for _, e := range entries {
		s := filepath.Join(src, e.Name())
		d := filepath.Join(dst, e.Name())
		info, err := os.Lstat(s)
		if err != nil {
			return err
		}
		if info.IsDir() {
			if err := copyDir(s, d); err != nil {
				return err
			}
		} else {
			if err := copyFile(s, d, info.Mode()); err != nil {
				return err
			}
		}
	}
	return nil
}
