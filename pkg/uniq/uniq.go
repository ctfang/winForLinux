package uniq

import (
	"bufio"
	"fmt"
	"github.com/user/win-tools/shared/bom"
	"io"
)

// Options for uniq
type Options struct {
	Count      bool // -c
	Repeated   bool // -d only duplicates
	UniqueOnly bool // -u only unique
}

// Uniq filters adjacent duplicate lines.
func Uniq(reader io.Reader, opts Options, writer io.Writer) error {
	scanner := bufio.NewScanner(bom.Reader(reader))
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 10*1024*1024)

	var prev string
	count := 0
	first := true

	flush := func() {
		if first {
			return
		}
		if opts.Repeated && count < 2 {
			return
		}
		if opts.UniqueOnly && count != 1 {
			return
		}
		if opts.Count {
			fmt.Fprintf(writer, "%7d %s\n", count, prev)
		} else {
			fmt.Fprintln(writer, prev)
		}
	}

	for scanner.Scan() {
		line := scanner.Text()
		if first {
			prev = line
			count = 1
			first = false
			continue
		}
		if line == prev {
			count++
			continue
		}
		flush()
		prev = line
		count = 1
	}
	flush()
	return scanner.Err()
}
