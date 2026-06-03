package argparse

import "testing"

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantOpts map[string]string
		wantPos  []string
	}{
		{
			name:     "short flag",
			args:     []string{"-l", "-a"},
			wantOpts: map[string]string{"l": "", "a": ""},
			wantPos:  nil,
		},
		{
			name:     "long flag",
			args:     []string{"--all", "--color=auto"},
			wantOpts: map[string]string{"all": "", "color": "auto"},
			wantPos:  nil,
		},
		{
			name:     "positional args",
			args:     []string{"file1.txt", "file2.txt"},
			wantOpts: map[string]string{},
			wantPos:  []string{"file1.txt", "file2.txt"},
		},
		{
			name:     "mixed",
			args:     []string{"-l", "dir1", "--color", "dir2"},
			wantOpts: map[string]string{"l": "", "color": ""},
			wantPos:  []string{"dir1", "dir2"},
		},
		{
			name:     "windows style",
			args:     []string{"/W", "/B"},
			wantOpts: map[string]string{"W": "", "B": ""},
			wantPos:  nil,
		},
		{
			name:     "combined short flags",
			args:     []string{"-lah"},
			wantOpts: map[string]string{"l": "", "a": "", "h": ""},
			wantPos:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts, pos := ParseArgs(tt.args)
			for k, v := range tt.wantOpts {
				if opts[k] != v {
					t.Errorf("opts[%q] = %q, want %q", k, opts[k], v)
				}
			}
			if len(pos) != len(tt.wantPos) {
				t.Errorf("positional args = %v, want %v", pos, tt.wantPos)
			}
		})
	}
}

func TestHasLinuxOption(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want bool
	}{
		{"no linux opts", []string{"/W", "/B"}, false},
		{"has linux opt", []string{"-l", "/W"}, true},
		{"all linux", []string{"-l", "-a", "-h"}, true},
		{"no opts", []string{"file.txt"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasLinuxOption(tt.args); got != tt.want {
				t.Errorf("HasLinuxOption(%v) = %v, want %v", tt.args, got, tt.want)
			}
		})
	}
}
