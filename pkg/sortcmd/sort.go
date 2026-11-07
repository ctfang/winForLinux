package sortcmd

import (
	"bufio"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
)

// Options for sort
type Options struct {
	Reverse  bool // -r
	Numeric  bool // -n
	Unique   bool // -u
	KeyField int  // -k N (1-based), 0 = whole line
}

// SortLines sorts lines from reader and writes to writer.
func SortLines(reader io.Reader, opts Options, writer io.Writer) error {
	scanner := bufio.NewScanner(reader)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 10*1024*1024)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	sort.SliceStable(lines, func(i, j int) bool {
		a, b := lines[i], lines[j]
		if opts.KeyField > 0 {
			a = fieldAt(a, opts.KeyField)
			b = fieldAt(b, opts.KeyField)
		}
		var less bool
		if opts.Numeric {
			na, ea := strconv.ParseFloat(strings.TrimSpace(a), 64)
			nb, eb := strconv.ParseFloat(strings.TrimSpace(b), 64)
			if ea == nil && eb == nil {
				less = na < nb
			} else {
				less = a < b
			}
		} else {
			less = a < b
		}
		if opts.Reverse {
			return !less && a != b
		}
		return less
	})

	if opts.Unique {
		var uniq []string
		var prev string
		first := true
		for _, line := range lines {
			key := line
			if opts.KeyField > 0 {
				key = fieldAt(line, opts.KeyField)
			}
			if first || key != prev {
				uniq = append(uniq, line)
				prev = key
				first = false
			}
		}
		lines = uniq
	}

	for _, line := range lines {
		fmt.Fprintln(writer, line)
	}
	return nil
}

func fieldAt(line string, n int) string {
	fields := strings.Fields(line)
	if n <= 0 || n > len(fields) {
		return ""
	}
	return fields[n-1]
}
