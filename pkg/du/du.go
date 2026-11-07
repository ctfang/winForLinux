package du

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Options for du
type Options struct {
	Human     bool // -h
	Summarize bool // -s
	All       bool // -a
	MaxDepth  int  // -d N; use MaxDepthSet
	DepthSet  bool
}

// Du reports disk usage for paths.
func Du(paths []string, opts Options, w io.Writer) error {
	if len(paths) == 0 {
		paths = []string{"."}
	}
	var firstErr error
	for _, p := range paths {
		if _, err := report(p, p, 0, opts, w); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

func report(root, path string, depth int, opts Options, w io.Writer) (int64, error) {
	info, err := os.Lstat(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du: %s: %v\n", path, err)
		return 0, err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return 0, nil
	}
	if !info.IsDir() {
		total := info.Size()
		if opts.All || path == root {
			printSize(total, path, opts, w)
		}
		return total, nil
	}

	var total int64
	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du: %s: %v\n", path, err)
		return 0, err
	}
	for _, e := range entries {
		child := filepath.Join(path, e.Name())
		sz, _ := report(root, child, depth+1, opts, w)
		total += sz
	}

	printThis := false
	if opts.Summarize {
		printThis = path == root
	} else if !opts.DepthSet || depth <= opts.MaxDepth {
		printThis = true
	}
	if printThis {
		printSize(total, path, opts, w)
	}
	return total, nil
}

func printSize(n int64, path string, opts Options, w io.Writer) {
	if opts.Human {
		fmt.Fprintf(w, "%s\t%s\n", humanSize(n), path)
		return
	}
	blocks := (n + 1023) / 1024
	fmt.Fprintf(w, "%d\t%s\n", blocks, path)
}

func humanSize(n int64) string {
	units := []string{"", "K", "M", "G", "T", "P"}
	v := float64(n)
	i := 0
	for v >= 1024 && i < len(units)-1 {
		v /= 1024
		i++
	}
	if i == 0 {
		return fmt.Sprintf("%d", n)
	}
	if v >= 10 {
		return fmt.Sprintf("%.0f%s", v, units[i])
	}
	return fmt.Sprintf("%.1f%s", v, units[i])
}