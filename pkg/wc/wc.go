package wc

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"
)

// Options for wc
type Options struct {
	Lines bool // -l
	Words bool // -w
	Bytes bool // -c
	Chars bool // -m
}

// Result of counting
type Result struct {
	Lines int
	Words int
	Bytes int
	Chars int
}

// Count counts lines/words/bytes/chars from reader.
func Count(reader io.Reader) (Result, error) {
	var r Result
	scanner := bufio.NewScanner(reader)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 10*1024*1024)
	for scanner.Scan() {
		line := scanner.Text()
		r.Lines++
		r.Words += len(strings.Fields(line))
		// scanner strips newline; approximate bytes with +1 for \n
		b := len(line) + 1
		r.Bytes += b
		r.Chars += utf8.RuneCountInString(line) + 1
	}
	if err := scanner.Err(); err != nil {
		return r, err
	}
	return r, nil
}

// CountBytes counts from raw bytes (more accurate for -c).
func CountBytes(data []byte) Result {
	var r Result
	r.Bytes = len(data)
	r.Chars = utf8.RuneCount(data)
	r.Lines = strings.Count(string(data), "\n")
	r.Words = len(strings.Fields(string(data)))
	return r
}

// Format formats result according to options. If no flag set, print l w c.
func Format(r Result, opts Options, name string) string {
	showAll := !opts.Lines && !opts.Words && !opts.Bytes && !opts.Chars
	var parts []string
	if showAll || opts.Lines {
		parts = append(parts, fmt.Sprintf("%d", r.Lines))
	}
	if showAll || opts.Words {
		parts = append(parts, fmt.Sprintf("%d", r.Words))
	}
	if showAll || opts.Bytes {
		parts = append(parts, fmt.Sprintf("%d", r.Bytes))
	}
	if opts.Chars && !showAll {
		parts = append(parts, fmt.Sprintf("%d", r.Chars))
	}
	out := strings.Join(parts, " ")
	if name != "" {
		out += " " + name
	}
	return out
}
