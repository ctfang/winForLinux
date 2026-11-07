package diff

import (
	"bufio"
	"fmt"
	"github.com/user/win-tools/shared/bom"
	"io"
	"os"
)

// Options for diff
type Options struct {
	Unified bool // -u
	Brief   bool // -q
}

// Diff compares two files and writes the result. Returns true if different.
func Diff(aPath, bPath string, opts Options, w io.Writer) (bool, error) {
	aLines, err := readLines(aPath)
	if err != nil {
		return false, err
	}
	bLines, err := readLines(bPath)
	if err != nil {
		return false, err
	}
	if equal(aLines, bLines) {
		return false, nil
	}
	if opts.Brief {
		fmt.Fprintf(w, "Files %s and %s differ\n", aPath, bPath)
		return true, nil
	}
	if opts.Unified {
		writeUnified(aPath, bPath, aLines, bLines, w)
	} else {
		writeNormal(aLines, bLines, w)
	}
	return true, nil
}

func readLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("diff: %s: %w", path, err)
	}
	defer f.Close()
	var lines []string
	sc := bufio.NewScanner(bom.Reader(f))
	buf := make([]byte, 0, 64*1024)
	sc.Buffer(buf, 10*1024*1024)
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	return lines, sc.Err()
}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func writeUnified(aName, bName string, a, b []string, w io.Writer) {
	fmt.Fprintf(w, "--- %s\n+++ %s\n", aName, bName)
	// Simple whole-file hunk (good enough for basic use)
	fmt.Fprintf(w, "@@ -1,%d +1,%d @@\n", len(a), len(b))
	lcs := lcsTable(a, b)
	i, j := len(a), len(b)
	var ops []string
	for i > 0 || j > 0 {
		if i > 0 && j > 0 && a[i-1] == b[j-1] {
			ops = append(ops, " "+a[i-1])
			i--
			j--
		} else if j > 0 && (i == 0 || lcs[i][j-1] >= lcs[i-1][j]) {
			ops = append(ops, "+"+b[j-1])
			j--
		} else if i > 0 {
			ops = append(ops, "-"+a[i-1])
			i--
		}
	}
	for k := len(ops) - 1; k >= 0; k-- {
		fmt.Fprintln(w, ops[k])
	}
}

func writeNormal(a, b []string, w io.Writer) {
	// very simple: print lines only in a / only in b via LCS
	lcs := lcsTable(a, b)
	i, j := len(a), len(b)
	var dels []string
	var adds []string
	type change struct {
		kind string
		line string
	}
	var seq []change
	for i > 0 || j > 0 {
		if i > 0 && j > 0 && a[i-1] == b[j-1] {
			seq = append(seq, change{" ", a[i-1]})
			i--
			j--
		} else if j > 0 && (i == 0 || lcs[i][j-1] >= lcs[i-1][j]) {
			seq = append(seq, change{"+", b[j-1]})
			j--
		} else {
			seq = append(seq, change{"-", a[i-1]})
			i--
		}
	}
	_ = dels
	_ = adds
	for k := len(seq) - 1; k >= 0; k-- {
		c := seq[k]
		if c.kind == "-" {
			fmt.Fprintf(w, "< %s\n", c.line)
		} else if c.kind == "+" {
			fmt.Fprintf(w, "> %s\n", c.line)
		}
	}
}

func lcsTable(a, b []string) [][]int {
	m, n := len(a), len(b)
	t := make([][]int, m+1)
	for i := range t {
		t[i] = make([]int, n+1)
	}
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if a[i-1] == b[j-1] {
				t[i][j] = t[i-1][j-1] + 1
			} else if t[i-1][j] >= t[i][j-1] {
				t[i][j] = t[i-1][j]
			} else {
				t[i][j] = t[i][j-1]
			}
		}
	}
	return t
}

// DiffReaders compares two readers (for tests).
func DiffReaders(a, b io.Reader, opts Options, w io.Writer, aName, bName string) (bool, error) {
	al, err := readAllLines(a)
	if err != nil {
		return false, err
	}
	bl, err := readAllLines(b)
	if err != nil {
		return false, err
	}
	if equal(al, bl) {
		return false, nil
	}
	if opts.Brief {
		fmt.Fprintf(w, "Files %s and %s differ\n", aName, bName)
		return true, nil
	}
	if opts.Unified {
		writeUnified(aName, bName, al, bl, w)
	} else {
		writeNormal(al, bl, w)
	}
	return true, nil
}

func readAllLines(r io.Reader) ([]string, error) {
	var lines []string
	sc := bufio.NewScanner(bom.Reader(r))
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	return lines, sc.Err()
}
