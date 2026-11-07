package tail

import (
	"bufio"
	"fmt"
	"github.com/user/win-tools/shared/bom"
	"io"
	"os"
	"time"
)

// Options for tail
type Options struct {
	Lines  int  // -n, default 10
	Follow bool // -f
}

// Tail writes the last N lines from reader to writer.
func Tail(reader io.Reader, opts Options, writer io.Writer) error {
	if opts.Lines <= 0 {
		opts.Lines = 10
	}

	scanner := bufio.NewScanner(bom.Reader(reader))
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 10*1024*1024)

	ring := make([]string, 0, opts.Lines)
	for scanner.Scan() {
		line := scanner.Text()
		if len(ring) < opts.Lines {
			ring = append(ring, line)
		} else {
			ring = append(ring[1:], line)
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	for _, line := range ring {
		fmt.Fprintln(writer, line)
	}
	return nil
}

// FollowFile appends new lines as the file grows (basic -f).
func FollowFile(path string, writer io.Writer) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// Seek to end
	if _, err := f.Seek(0, io.SeekEnd); err != nil {
		return err
	}

	reader := bufio.NewReader(f)
	for {
		line, err := reader.ReadString('\n')
		if len(line) > 0 {
			fmt.Fprint(writer, line)
		}
		if err == io.EOF {
			time.Sleep(200 * time.Millisecond)
			continue
		}
		if err != nil {
			return err
		}
	}
}
