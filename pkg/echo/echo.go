package echo

import (
	"fmt"
	"io"
	"strings"
)

// Options for echo
type Options struct {
	NoNewline bool // -n
	EnableEsc bool // -e interpret backslash escapes
}

// Echo writes arguments joined by space to writer.
func Echo(args []string, opts Options, w io.Writer) error {
	s := strings.Join(args, " ")
	if opts.EnableEsc {
		s = interpretEscapes(s)
	}
	if opts.NoNewline {
		_, err := io.WriteString(w, s)
		return err
	}
	_, err := fmt.Fprintln(w, s)
	return err
}

func interpretEscapes(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for i := 0; i < len(s); i++ {
		if s[i] != '\\' || i+1 >= len(s) {
			b.WriteByte(s[i])
			continue
		}
		i++
		switch s[i] {
		case 'n':
			b.WriteByte('\n')
		case 't':
			b.WriteByte('\t')
		case 'r':
			b.WriteByte('\r')
		case '\\':
			b.WriteByte('\\')
		case 'a':
			b.WriteByte('\a')
		case 'b':
			b.WriteByte('\b')
		case '0':
			b.WriteByte(0)
		case 'c':
			return b.String()
		default:
			b.WriteByte('\\')
			b.WriteByte(s[i])
		}
	}
	return b.String()
}