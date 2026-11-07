package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/user/win-tools/pkg/wc"
)

var version = "dev"

func expandShort(args []string) []string {
	var out []string
	for _, arg := range args {
		if strings.HasPrefix(arg, "--") {
			out = append(out, arg)
			continue
		}
		if strings.HasPrefix(arg, "-") && len(arg) > 2 {
			all := true
			for _, ch := range arg[1:] {
				if (ch < 'a' || ch > 'z') && (ch < 'A' || ch > 'Z') {
					all = false
					break
				}
			}
			if all {
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

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-V") {
		fmt.Println("win-tools wc", version)
		os.Exit(0)
	}

	opts := wc.Options{}
	var files []string
	for _, arg := range expandShort(os.Args[1:]) {
		switch arg {
		case "-l", "--lines":
			opts.Lines = true
		case "-w", "--words":
			opts.Words = true
		case "-c", "--bytes":
			opts.Bytes = true
		case "-m", "--chars":
			opts.Chars = true
		case "-h", "--help":
			fmt.Fprintf(os.Stderr, "Usage: wc [-lwc] [FILE...]\n")
			os.Exit(0)
		default:
			files = append(files, arg)
		}
	}

	if len(files) == 0 {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "wc: %v\n", err)
			os.Exit(1)
		}
		r := wc.CountBytes(data)
		fmt.Println(wc.Format(r, opts, ""))
		return
	}

	var total wc.Result
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "wc: %s: %v\n", f, err)
			continue
		}
		r := wc.CountBytes(data)
		total.Lines += r.Lines
		total.Words += r.Words
		total.Bytes += r.Bytes
		total.Chars += r.Chars
		name := f
		if len(files) == 1 {
			name = f
		}
		fmt.Println(wc.Format(r, opts, name))
	}
	if len(files) > 1 {
		fmt.Println(wc.Format(total, opts, "total"))
	}
}
