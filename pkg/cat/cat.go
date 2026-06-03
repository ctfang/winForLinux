package cat

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Options cat command options
type Options struct {
	Number       bool // -n
	ShowEnds     bool // -E
	ShowTabs     bool // -T
	SqueezeBlank bool // -s
}

// Cat reads from reader and outputs to writer
func Cat(reader io.Reader, opts Options, writer io.Writer) error {
	scanner := bufio.NewScanner(reader)
	lineNum := 0
	prevBlank := false

	for scanner.Scan() {
		line := scanner.Text()

		// Squeeze blank lines
		if opts.SqueezeBlank {
			isBlank := line == ""
			if isBlank && prevBlank {
				continue
			}
			prevBlank = isBlank
		}

		lineNum++

		// Line numbers
		if opts.Number {
			fmt.Fprintf(writer, "%6d\t", lineNum)
		}

		// Show tabs
		if opts.ShowTabs {
			line = strings.ReplaceAll(line, "\t", "^I")
		}

		fmt.Fprint(writer, line)

		// Show line ends
		if opts.ShowEnds {
			fmt.Fprint(writer, "$")
		}

		fmt.Fprintln(writer)
	}

	return scanner.Err()
}
