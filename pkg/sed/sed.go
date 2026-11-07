package sed

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Options for sed
type Options struct {
	Expressions []string
	InPlace     bool
	Quiet       bool // -n
}

type cmd struct {
	Op      byte
	Re      *regexp.Regexp
	Repl    string
	Global  bool
	AddrNum int
	AddrRe  *regexp.Regexp
	AddrAll bool
}

// Run applies sed scripts to reader, writing to writer.
func Run(reader io.Reader, opts Options, writer io.Writer) error {
	cmds, err := parseScripts(opts.Expressions)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(reader)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 10*1024*1024)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		deleted := false
		for _, c := range cmds {
			if !addrMatch(c, lineNum, line) {
				continue
			}
			switch c.Op {
			case 's':
				if c.Global {
					line = c.Re.ReplaceAllString(line, c.Repl)
				} else {
					line = replaceFirst(line, c.Re, c.Repl)
				}
			case 'd':
				deleted = true
			case 'p':
				fmt.Fprintln(writer, line)
			case 'q':
				if !opts.Quiet && !deleted {
					fmt.Fprintln(writer, line)
				}
				return nil
			}
			if deleted {
				break
			}
		}
		if deleted {
			continue
		}
		if !opts.Quiet {
			fmt.Fprintln(writer, line)
		}
	}
	return scanner.Err()
}

func replaceFirst(s string, re *regexp.Regexp, repl string) string {
	loc := re.FindStringIndex(s)
	if loc == nil {
		return s
	}
	replaced := re.ReplaceAllString(s[loc[0]:loc[1]], repl)
	return s[:loc[0]] + replaced + s[loc[1]:]
}

func addrMatch(c cmd, lineNum int, line string) bool {
	if c.AddrAll {
		return true
	}
	if c.AddrNum > 0 {
		return lineNum == c.AddrNum
	}
	if c.AddrRe != nil {
		return c.AddrRe.MatchString(line)
	}
	return true
}

func parseScripts(exprs []string) ([]cmd, error) {
	var out []cmd
	for _, e := range exprs {
		c, err := parseOne(strings.TrimSpace(e))
		if err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("sed: no expression")
	}
	return out, nil
}

func parseOne(e string) (cmd, error) {
	var c cmd
	c.AddrAll = true
	rest := e

	if len(rest) > 0 && rest[0] >= '0' && rest[0] <= '9' {
		j := 0
		for j < len(rest) && rest[j] >= '0' && rest[j] <= '9' {
			j++
		}
		n, _ := strconv.Atoi(rest[:j])
		c.AddrNum = n
		c.AddrAll = false
		rest = rest[j:]
	} else if len(rest) > 0 && rest[0] == '/' {
		end := strings.IndexByte(rest[1:], '/')
		if end < 0 {
			return c, fmt.Errorf("sed: unterminated address")
		}
		pat := rest[1 : 1+end]
		re, err := regexp.Compile(pat)
		if err != nil {
			return c, fmt.Errorf("sed: %w", err)
		}
		c.AddrRe = re
		c.AddrAll = false
		rest = rest[2+end:]
	}

	if rest == "" {
		return c, fmt.Errorf("sed: missing command")
	}
	c.Op = rest[0]
	rest = rest[1:]

	switch c.Op {
	case 's':
		if rest == "" {
			return c, fmt.Errorf("sed: bad substitute")
		}
		sep := rest[0]
		parts := splitSep(rest[1:], sep)
		if len(parts) < 2 {
			return c, fmt.Errorf("sed: bad substitute")
		}
		pat, repl := parts[0], parts[1]
		flags := ""
		if len(parts) > 2 {
			flags = parts[2]
		}
		patUse := pat
		if strings.Contains(flags, "i") {
			patUse = "(?i)" + pat
		}
		re, err := regexp.Compile(patUse)
		if err != nil {
			return c, fmt.Errorf("sed: %w", err)
		}
		c.Re = re
		c.Repl = unescapeRepl(repl)
		c.Global = strings.Contains(flags, "g")
	case 'd', 'p', 'q':
		// ok
	default:
		return c, fmt.Errorf("sed: unsupported command: %c", c.Op)
	}
	return c, nil
}

func splitSep(s string, sep byte) []string {
	var parts []string
	var cur strings.Builder
	esc := false
	for i := 0; i < len(s); i++ {
		if esc {
			cur.WriteByte(s[i])
			esc = false
			continue
		}
		if s[i] == '\\' {
			esc = true
			cur.WriteByte(s[i])
			continue
		}
		if s[i] == sep {
			parts = append(parts, cur.String())
			cur.Reset()
			continue
		}
		cur.WriteByte(s[i])
	}
	parts = append(parts, cur.String())
	return parts
}

func unescapeRepl(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' && i+1 < len(s) {
			i++
			switch s[i] {
			case 'n':
				b.WriteByte('\n')
			case 't':
				b.WriteByte('\t')
			case '\\':
				b.WriteByte('\\')
			case '&':
				b.WriteByte('&')
			case '1', '2', '3', '4', '5', '6', '7', '8', '9':
				b.WriteByte('$')
				b.WriteByte(s[i])
			default:
				b.WriteByte(s[i])
			}
			continue
		}
		if s[i] == '&' {
			b.WriteString("$0")
			continue
		}
		b.WriteByte(s[i])
	}
	return b.String()
}

// RunFile runs sed on a file; if InPlace, rewrites the file.
func RunFile(path string, opts Options) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("sed: %s: %w", path, err)
	}
	var buf strings.Builder
	if err := Run(strings.NewReader(string(data)), opts, &buf); err != nil {
		return err
	}
	if opts.InPlace {
		return os.WriteFile(path, []byte(buf.String()), 0644)
	}
	_, err = io.WriteString(os.Stdout, buf.String())
	return err
}