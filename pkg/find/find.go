package find

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Options for find
type Options struct {
	Name     string
	IName    string
	Type     string
	MaxDepth int  // -1 = unlimited
	MinDepth int
	Size     string
	MTime    string
}

// Find walks paths and prints matches to writer.
func Find(paths []string, opts Options, w io.Writer) error {
	if len(paths) == 0 {
		paths = []string{"."}
	}
	var firstErr error
	for _, root := range paths {
		if err := walk(root, opts, w); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

func walk(root string, opts Options, w io.Writer) error {
	root = filepath.Clean(root)
	info, err := os.Lstat(root)
	if err != nil {
		return fmt.Errorf("find: '%s': %w", root, err)
	}
	return walkAt(root, info, 0, opts, w)
}

func walkAt(path string, info os.FileInfo, depth int, opts Options, w io.Writer) error {
	if opts.MaxDepth >= 0 && depth > opts.MaxDepth {
		return nil
	}
	if depth >= opts.MinDepth {
		if match(path, info, opts) {
			fmt.Fprintln(w, path)
		}
	}
	if !info.IsDir() {
		return nil
	}
	if opts.MaxDepth >= 0 && depth >= opts.MaxDepth {
		return nil
	}
	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "find: '%s': %v\n", path, err)
		return nil
	}
	for _, e := range entries {
		child := filepath.Join(path, e.Name())
		ci, err := e.Info()
		if err != nil {
			continue
		}
		if err := walkAt(child, ci, depth+1, opts, w); err != nil {
			return err
		}
	}
	return nil
}

func match(path string, info os.FileInfo, opts Options) bool {
	base := filepath.Base(path)
	if opts.Name != "" {
		ok, _ := filepath.Match(opts.Name, base)
		if !ok {
			return false
		}
	}
	if opts.IName != "" {
		ok, _ := filepath.Match(strings.ToLower(opts.IName), strings.ToLower(base))
		if !ok {
			return false
		}
	}
	if opts.Type != "" {
		switch opts.Type {
		case "f":
			if info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
				return false
			}
		case "d":
			if !info.IsDir() {
				return false
			}
		case "l":
			if info.Mode()&os.ModeSymlink == 0 {
				return false
			}
		}
	}
	if opts.Size != "" && !matchSize(info.Size(), opts.Size) {
		return false
	}
	if opts.MTime != "" && !matchMTime(info.ModTime(), opts.MTime) {
		return false
	}
	return true
}

func matchSize(size int64, expr string) bool {
	if len(expr) == 0 {
		return true
	}
	cmp := byte(0)
	s := expr
	if s[0] == '+' || s[0] == '-' {
		cmp = s[0]
		s = s[1:]
	}
	if len(s) == 0 {
		return false
	}
	unit := s[len(s)-1]
	numStr := s
	mult := int64(512)
	switch unit {
	case 'c':
		numStr = s[:len(s)-1]
		mult = 1
	case 'k', 'K':
		numStr = s[:len(s)-1]
		mult = 1024
	case 'M':
		numStr = s[:len(s)-1]
		mult = 1024 * 1024
	case 'G':
		numStr = s[:len(s)-1]
		mult = 1024 * 1024 * 1024
	case 'b':
		numStr = s[:len(s)-1]
		mult = 512
	default:
		if unit < '0' || unit > '9' {
			return false
		}
	}
	n, err := strconv.ParseInt(numStr, 10, 64)
	if err != nil {
		return false
	}
	target := n * mult
	switch cmp {
	case '+':
		return size > target
	case '-':
		return size < target
	default:
		if mult == 512 && (unit < 'A' || unit > 'z') {
			blocks := (size + 511) / 512
			return blocks == n
		}
		return size == target
	}
}

func matchMTime(mod time.Time, expr string) bool {
	cmp := byte(0)
	s := expr
	if len(s) > 0 && (s[0] == '+' || s[0] == '-') {
		cmp = s[0]
		s = s[1:]
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return false
	}
	days := int(time.Since(mod).Hours() / 24)
	switch cmp {
	case '+':
		return days > n
	case '-':
		return days < n
	default:
		return days == n
	}
}