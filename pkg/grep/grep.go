package grep

import (
	"bufio"
	"fmt"
	"io"
	"regexp"

	"github.com/user/win-tools/shared/color"
)

// Options grep command options
type Options struct {
	IgnoreCase  bool
	LineNumber  bool
	Count       bool
	Invert      bool
	OnlyMatch   bool
	WordRegexp  bool
	FixedString bool
	Color       string
}

// Grep reads from reader and matches pattern
func Grep(reader io.Reader, pattern string, opts Options, writer io.Writer) error {
	// Compile regex
	regexPattern := pattern
	if opts.FixedString {
		regexPattern = regexp.QuoteMeta(pattern)
	}
	if opts.WordRegexp {
		regexPattern = `\b` + regexPattern + `\b`
	}
	if opts.IgnoreCase {
		regexPattern = "(?i)" + regexPattern
	}

	re, err := regexp.Compile(regexPattern)
	if err != nil {
		return fmt.Errorf("invalid pattern: %v", err)
	}

	scanner := bufio.NewScanner(reader)
	lineNum := 0
	matchCount := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		matched := re.MatchString(line)

		if opts.Invert {
			matched = !matched
		}

		if matched {
			matchCount++

			if opts.Count {
				continue
			}

			if opts.OnlyMatch {
				matches := re.FindAllString(line, -1)
				for _, m := range matches {
					if opts.LineNumber {
						fmt.Fprintf(writer, "%d:%s\n", lineNum, m)
					} else {
						fmt.Fprintf(writer, "%s\n", m)
					}
				}
				continue
			}

			var output string
			if opts.LineNumber {
				output = fmt.Sprintf("%d:", lineNum)
			}

			// Highlight matches
			if opts.Color != "never" && opts.Color != "" {
				output += re.ReplaceAllString(line, color.Colorize("$0", color.Red))
			} else {
				output += line
			}

			fmt.Fprintln(writer, output)
		}
	}

	if opts.Count {
		fmt.Fprintf(writer, "%d\n", matchCount)
	}

	return scanner.Err()
}
