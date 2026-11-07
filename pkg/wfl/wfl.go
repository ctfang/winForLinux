// Package wfl implements a small Linux-style command runner for Windows.
// It supports |, &&, ||, ; and & between external commands (preferring
// binaries installed by win-tools when present on PATH).
package wfl

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

// Op represents a shell operator between commands.
type Op int

const (
	OpNone Op = iota
	OpPipe    // |
	OpAnd     // &&
	OpOr      // ||
	OpSeq     // ;
	OpBg      // &
)

func (o Op) String() string {
	switch o {
	case OpPipe:
		return "|"
	case OpAnd:
		return "&&"
	case OpOr:
		return "||"
	case OpSeq:
		return ";"
	case OpBg:
		return "&"
	default:
		return ""
	}
}

// CmdSpec is one command with argv.
type CmdSpec struct {
	Args []string
	Op   Op // operator AFTER this command (OpNone for last)
}

// Tokenize splits a command line into tokens, respecting single/double quotes.
func Tokenize(line string) ([]string, error) {
	var tokens []string
	var cur strings.Builder
	inSingle, inDouble := false, false
	escaped := false

	flush := func() {
		if cur.Len() > 0 {
			tokens = append(tokens, cur.String())
			cur.Reset()
		}
	}

	for i := 0; i < len(line); i++ {
		ch := line[i]
		if escaped {
			cur.WriteByte(ch)
			escaped = false
			continue
		}
		if ch == '\\' && !inSingle {
			// On Windows paths, '\' is a separator. Only treat it as an escape
			// when the next character is a shell metacharacter or quote.
			if i+1 < len(line) {
				next := line[i+1]
				if next == '"' || next == '\'' || next == '\\' ||
					next == '|' || next == '&' || next == ';' ||
					next == ' ' || next == '\t' {
					escaped = true
					continue
				}
			}
			// literal backslash (e.g. C:\Users\...)
		}
		if ch == '\'' && !inDouble {
			inSingle = !inSingle
			continue
		}
		if ch == '"' && !inSingle {
			inDouble = !inDouble
			continue
		}
		if !inSingle && !inDouble {
			if ch == '&' && i+1 < len(line) && line[i+1] == '&' {
				flush()
				tokens = append(tokens, "&&")
				i++
				continue
			}
			if ch == '|' && i+1 < len(line) && line[i+1] == '|' {
				flush()
				tokens = append(tokens, "||")
				i++
				continue
			}
			if ch == '|' || ch == '&' || ch == ';' {
				flush()
				tokens = append(tokens, string(ch))
				continue
			}
			if ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' {
				flush()
				continue
			}
		}
		cur.WriteByte(ch)
	}
	if escaped {
		cur.WriteByte('\\')
	}
	if inSingle || inDouble {
		return nil, fmt.Errorf("wfl: unmatched quote")
	}
	flush()
	return tokens, nil
}

// Parse converts tokens into a list of CmdSpec with trailing operators.
func Parse(tokens []string) ([]CmdSpec, error) {
	var specs []CmdSpec
	var cur []string

	flush := func(op Op) error {
		if len(cur) == 0 {
			if op != OpNone {
				return fmt.Errorf("wfl: unexpected operator %s", op)
			}
			return nil
		}
		specs = append(specs, CmdSpec{Args: append([]string{}, cur...), Op: op})
		cur = nil
		return nil
	}

	for _, tok := range tokens {
		var op Op
		switch tok {
		case "|":
			op = OpPipe
		case "&&":
			op = OpAnd
		case "||":
			op = OpOr
		case ";":
			op = OpSeq
		case "&":
			op = OpBg
		default:
			cur = append(cur, tok)
			continue
		}
		if err := flush(op); err != nil {
			return nil, err
		}
	}
	if err := flush(OpNone); err != nil {
		return nil, err
	}
	if len(specs) == 0 {
		return nil, fmt.Errorf("wfl: empty command")
	}
	return specs, nil
}

func lookPath(name string) (string, error) {
	if strings.ContainsAny(name, `/\`) {
		return name, nil
	}
	// Prefer sibling win-tools binaries next to wfl.exe
	if self, err := os.Executable(); err == nil {
		dir := filepath.Dir(self)
		candidates := []string{filepath.Join(dir, name+".exe"), filepath.Join(dir, name)}
		for _, c := range candidates {
			if st, err := os.Stat(c); err == nil && !st.IsDir() {
				return c, nil
			}
		}
	}
	if runtime.GOOS == "windows" {
		if p, err := exec.LookPath(name + ".exe"); err == nil {
			return p, nil
		}
	}
	return exec.LookPath(name)
}

func isBuiltin(name string) bool {
	switch name {
	case "echo", "true", "false", "cd":
		return true
	}
	return false
}

func runBuiltin(args []string, stdout, stderr io.Writer) int {
	switch args[0] {
	case "echo":
		fmt.Fprintln(stdout, strings.Join(args[1:], " "))
		return 0
	case "true":
		return 0
	case "false":
		return 1
	case "cd":
		dir := ""
		if len(args) > 1 {
			dir = args[1]
		} else {
			dir, _ = os.UserHomeDir()
		}
		if err := os.Chdir(dir); err != nil {
			fmt.Fprintf(stderr, "wfl: cd: %v\n", err)
			return 1
		}
		return 0
	}
	return 127
}

// Run executes the parsed command chain with Linux-like operator semantics.
func Run(specs []CmdSpec, stdin io.Reader, stdout, stderr io.Writer) int {
	exitCode := 0
	i := 0
	for i < len(specs) {
		// Gather pipe chain: consecutive OpPipe links
		j := i
		for j < len(specs)-1 && specs[j].Op == OpPipe {
			j++
		}
		chain := specs[i : j+1]
		bg := specs[j].Op == OpBg

		if bg {
			chainCopy := append([]CmdSpec{}, chain...)
			go func() {
				_ = runPipeline(chainCopy, stdin, stdout, stderr)
			}()
			fmt.Fprintf(stderr, "[bg] started\n")
			exitCode = 0
		} else {
			exitCode = runPipeline(chain, stdin, stdout, stderr)
		}

		op := specs[j].Op
		next := j + 1
		switch op {
		case OpAnd:
			if exitCode != 0 {
				// skip subsequent &&-linked segments
				for next < len(specs) && specs[next-1].Op == OpAnd {
					// advance past next pipe-chain under &&
					k := next
					for k < len(specs)-1 && specs[k].Op == OpPipe {
						k++
					}
					next = k + 1
				}
			}
		case OpOr:
			if exitCode == 0 {
				for next < len(specs) && specs[next-1].Op == OpOr {
					k := next
					for k < len(specs)-1 && specs[k].Op == OpPipe {
						k++
					}
					next = k + 1
				}
			}
		}
		i = next
	}
	return exitCode
}

func runPipeline(chain []CmdSpec, stdin io.Reader, stdout, stderr io.Writer) int {
	n := len(chain)
	if n == 0 {
		return 0
	}

	// Single builtin without pipe
	if n == 1 && isBuiltin(chain[0].Args[0]) {
		return runBuiltin(chain[0].Args, stdout, stderr)
	}

	cmds := make([]*exec.Cmd, n)
	var closers []io.Closer

	var prevOut io.Reader = stdin
	for idx, spec := range chain {
		isLast := idx == n-1

		if isBuiltin(spec.Args[0]) {
			// Builtin in a pipeline: materialize via pipe goroutine
			pr, pw := io.Pipe()
			args := append([]string{}, spec.Args...)
			in := prevOut
			go func() {
				defer pw.Close()
				if in != nil && in != os.Stdin {
					_, _ = io.Copy(io.Discard, in) // drain unused stdin for builtins except echo
				}
				_ = runBuiltin(args, pw, stderr)
			}()
			if isLast {
				_, _ = io.Copy(stdout, pr)
				pr.Close()
				return 0
			}
			prevOut = pr
			closers = append(closers, pr)
			continue
		}

		path, err := lookPath(spec.Args[0])
		if err != nil {
			fmt.Fprintf(stderr, "wfl: %s: command not found\n", spec.Args[0])
			return 127
		}
		cmd := exec.Command(path, spec.Args[1:]...)
		cmd.Stdin = prevOut
		cmd.Stderr = stderr
		if isLast {
			cmd.Stdout = stdout
		} else {
			stdoutPipe, err := cmd.StdoutPipe()
			if err != nil {
				fmt.Fprintln(stderr, err)
				return 1
			}
			prevOut = stdoutPipe
			closers = append(closers, stdoutPipe)
		}
		cmds[idx] = cmd
	}

	// Start all non-nil commands
	for _, cmd := range cmds {
		if cmd == nil {
			continue
		}
		if err := cmd.Start(); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
	}

	code := 0
	var wg sync.WaitGroup
	for idx, cmd := range cmds {
		if cmd == nil {
			continue
		}
		wg.Add(1)
		go func(i int, c *exec.Cmd) {
			defer wg.Done()
			err := c.Wait()
			if i == n-1 {
				if err != nil {
					if ee, ok := err.(*exec.ExitError); ok {
						code = ee.ExitCode()
					} else {
						code = 1
					}
				}
			}
		}(idx, cmd)
	}
	wg.Wait()
	for _, c := range closers {
		_ = c.Close()
	}
	return code
}

// RunLine tokenizes, parses and runs a full command line string.
func RunLine(line string, stdin io.Reader, stdout, stderr io.Writer) int {
	line = strings.TrimSpace(line)
	if line == "" {
		return 0
	}
	tokens, err := Tokenize(line)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	specs, err := Parse(tokens)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	return Run(specs, stdin, stdout, stderr)
}
