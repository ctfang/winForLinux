package xargs

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

// Options for xargs
type Options struct {
	MaxArgs      int    // -n
	Replace      string // -I {}
	NullSep      bool   // -0
	Delimiter    string // -d
	NoRunIfEmpty bool   // -r
}

// Run reads items from reader and executes command with them.
// command[0] is the program; remaining are fixed args before items.
func Run(reader io.Reader, command []string, opts Options, stdout, stderr io.Writer) error {
	if len(command) == 0 {
		command = []string{"echo"}
	}
	items, err := readItems(reader, opts)
	if err != nil {
		return err
	}
	if len(items) == 0 {
		if opts.NoRunIfEmpty {
			return nil
		}
		if opts.Replace != "" {
			return nil
		}
		return runOnce(command, nil, stdout, stderr)
	}

	if opts.Replace != "" {
		for _, item := range items {
			args := replaceArgs(command, opts.Replace, item)
			if err := runOnce(args, nil, stdout, stderr); err != nil {
				return err
			}
		}
		return nil
	}

	maxArgs := opts.MaxArgs
	if maxArgs <= 0 {
		maxArgs = 64
	}
	for i := 0; i < len(items); i += maxArgs {
		end := i + maxArgs
		if end > len(items) {
			end = len(items)
		}
		batch := items[i:end]
		if err := runOnce(command, batch, stdout, stderr); err != nil {
			return err
		}
	}
	return nil
}

func replaceArgs(command []string, token, item string) []string {
	out := make([]string, len(command))
	for i, a := range command {
		out[i] = strings.ReplaceAll(a, token, item)
	}
	return out
}

func runOnce(command []string, extra []string, stdout, stderr io.Writer) error {
	args := append(append([]string{}, command...), extra...)
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.Stdin = nil
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("xargs: %w", err)
	}
	return nil
}

func readItems(r io.Reader, opts Options) ([]string, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, nil
	}

	if opts.NullSep {
		parts := bytes.Split(data, []byte{0})
		var items []string
		for _, p := range parts {
			if len(p) == 0 {
				continue
			}
			items = append(items, string(p))
		}
		return items, nil
	}

	if opts.Delimiter != "" {
		sep := opts.Delimiter
		switch sep {
		case "\\n":
			sep = "\n"
		case "\\t":
			sep = "\t"
		case "\\0":
			sep = "\x00"
		}
		parts := strings.Split(string(data), sep)
		var items []string
		for _, p := range parts {
			p = strings.TrimRight(p, "\r\n")
			if p == "" {
				continue
			}
			items = append(items, p)
		}
		return items, nil
	}

	var items []string
	sc := bufio.NewScanner(bytes.NewReader(data))
	sc.Split(bufio.ScanWords)
	for sc.Scan() {
		items = append(items, sc.Text())
	}
	return items, sc.Err()
}