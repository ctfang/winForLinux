package grep

import (
	"bytes"
	"strings"
	"testing"
)

func TestGrepLines(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		pattern string
		opts    Options
		want    []string
	}{
		{
			name:    "basic match",
			input:   "hello\nworld\nhello world",
			pattern: "hello",
			opts:    Options{},
			want:    []string{"hello", "hello world"},
		},
		{
			name:    "case insensitive",
			input:   "Hello\nworld\nHELLO",
			pattern: "hello",
			opts:    Options{IgnoreCase: true},
			want:    []string{"Hello", "HELLO"},
		},
		{
			name:    "with line numbers",
			input:   "line1\nline2\nline3",
			pattern: "line2",
			opts:    Options{LineNumber: true},
			want:    []string{"2:line2"},
		},
		{
			name:    "invert match",
			input:   "apple\nbanana\ncherry",
			pattern: "banana",
			opts:    Options{Invert: true},
			want:    []string{"apple", "cherry"},
		},
		{
			name:    "count only",
			input:   "a\nb\nc\na",
			pattern: "a",
			opts:    Options{Count: true},
			want:    []string{"2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := Grep(strings.NewReader(tt.input), tt.pattern, tt.opts, &buf)
			if err != nil {
				t.Fatalf("Grep() error = %v", err)
			}
			lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
			if len(lines) != len(tt.want) {
				t.Fatalf("got %d lines, want %d", len(lines), len(tt.want))
			}
			for i, line := range lines {
				if line != tt.want[i] {
					t.Errorf("line %d = %q, want %q", i, line, tt.want[i])
				}
			}
		})
	}
}
