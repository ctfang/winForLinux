package grep

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

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
	Recursive   bool
	Color       string
}

// Grep reads from reader and matches pattern. Returns whether any match was found.
func Grep(reader io.Reader, pattern string, opts Options, writer io.Writer) (bool, error) {
	return GrepNamed(reader, "", pattern, opts, writer)
}

// GrepNamed is like Grep but prefixes matches with filename when name is non-empty
// (used for multi-file / recursive searches).
func GrepNamed(reader io.Reader, name string, pattern string, opts Options, writer io.Writer) (bool, error) {
	re, err := compilePattern(pattern, opts)
	if err != nil {
		return false, err
	}

	scanner := bufio.NewScanner(reader)
	// Allow long lines
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 10*1024*1024)

	lineNum := 0
	matchCount := 0
	found := false

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		matched := re.MatchString(line)

		if opts.Invert {
			matched = !matched
		}

		if matched {
			found = true
			matchCount++

			if opts.Count {
				continue
			}

			prefix := ""
			if name != "" {
				prefix = name + ":"
			}

			if opts.OnlyMatch {
				matches := re.FindAllString(line, -1)
				for _, m := range matches {
					if opts.LineNumber {
						fmt.Fprintf(writer, "%s%d:%s\n", prefix, lineNum, m)
					} else {
						fmt.Fprintf(writer, "%s%s\n", prefix, m)
					}
				}
				continue
			}

			var output string
			if opts.LineNumber {
				output = fmt.Sprintf("%s%d:", prefix, lineNum)
			} else {
				output = prefix
			}

			if opts.Color == "always" {
				output += re.ReplaceAllString(line, color.Colorize("$0", color.Red))
			} else {
				output += line
			}

			fmt.Fprintln(writer, output)
		}
	}

	if opts.Count {
		if name != "" {
			fmt.Fprintf(writer, "%s:%d\n", name, matchCount)
		} else {
			fmt.Fprintf(writer, "%d\n", matchCount)
		}
	}

	return found || (opts.Count && matchCount > 0), scanner.Err()
}

func compilePattern(pattern string, opts Options) (*regexp.Regexp, error) {
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
		return nil, fmt.Errorf("invalid pattern: %v", err)
	}
	return re, nil
}

// GrepRecursive walks root and greps regular files.
func GrepRecursive(root string, pattern string, opts Options, writer io.Writer) (bool, error) {
	foundAny := false
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Fprintf(os.Stderr, "grep: %s: %v\n", path, err)
			return nil
		}
		if info.IsDir() {
			base := info.Name()
			if base == ".git" || base == "node_modules" || base == "vendor" {
				return filepath.SkipDir
			}
			return nil
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "grep: %s: %v\n", path, err)
			return nil
		}
		defer f.Close()
		display := path
		found, gerr := GrepNamed(f, display, pattern, opts, writer)
		if gerr != nil {
			fmt.Fprintf(os.Stderr, "grep: %s: %v\n", path, gerr)
			return nil
		}
		if found {
			foundAny = true
		}
		return nil
	})
	return foundAny, err
}

// ExpandShortFlags expands combined short options like -in into -i -n.
// Stops treating as flags once a non-flag positional token is needed as pattern.
func ExpandShortFlags(args []string) []string {
	var out []string
	for _, arg := range args {
		if strings.HasPrefix(arg, "--") {
			out = append(out, arg)
			continue
		}
		if strings.HasPrefix(arg, "-") && len(arg) > 2 && !strings.Contains(arg[1:], "-") {
			// Combined short flags: -in -> -i -n (skip digits-only like -10 for other cmds)
			allFlags := true
			for _, ch := range arg[1:] {
				if (ch < 'a' || ch > 'z') && (ch < 'A' || ch > 'Z') {
					allFlags = false
					break
				}
			}
			if allFlags {
				for _, ch := range arg[1:] {
					out = append(out, "-"+string(ch))
				}
				continue
			}
		}
		out = append(out, arg)
	}
	return out
}
