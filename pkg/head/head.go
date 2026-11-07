package head

import (
	"bufio"
	"fmt"
	"io"
)

// Options for head
type Options struct {
	Lines int // -n, default 10
	Bytes int // -c, 0 means unused
}

// Head writes the first N lines or bytes from reader to writer.
func Head(reader io.Reader, opts Options, writer io.Writer) error {
	if opts.Lines <= 0 && opts.Bytes <= 0 {
		opts.Lines = 10
	}

	if opts.Bytes > 0 {
		buf := make([]byte, opts.Bytes)
		n, err := io.ReadFull(reader, buf)
		if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
			return err
		}
		_, werr := writer.Write(buf[:n])
		return werr
	}

	scanner := bufio.NewScanner(reader)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 10*1024*1024)
	count := 0
	for scanner.Scan() {
		count++
		fmt.Fprintln(writer, scanner.Text())
		if count >= opts.Lines {
			break
		}
	}
	return scanner.Err()
}
