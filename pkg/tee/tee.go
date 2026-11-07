package tee

import (
	"fmt"
	"io"
	"os"
)

// Options for tee
type Options struct {
	Append bool // -a
}

// Tee copies reader to writer and to each file path.
func Tee(reader io.Reader, paths []string, opts Options, writer io.Writer) error {
	var files []*os.File
	defer func() {
		for _, f := range files {
			f.Close()
		}
	}()

	flag := os.O_CREATE | os.O_WRONLY
	if opts.Append {
		flag |= os.O_APPEND
	} else {
		flag |= os.O_TRUNC
	}

	writers := []io.Writer{writer}
	for _, p := range paths {
		f, err := os.OpenFile(p, flag, 0644)
		if err != nil {
			return fmt.Errorf("tee: %s: %w", p, err)
		}
		files = append(files, f)
		writers = append(writers, f)
	}

	mw := io.MultiWriter(writers...)
	_, err := io.Copy(mw, reader)
	return err
}