package ls

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseOptions(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want Options
	}{
		{
			name: "windows mode",
			args: []string{"/W", "/B"},
			want: Options{WindowsMode: true, Wide: true, Bare: true},
		},
		{
			name: "linux mode",
			args: []string{"-l", "-a"},
			want: Options{Long: true, All: true},
		},
		{
			name: "no args",
			args: []string{},
			want: Options{WindowsMode: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseOptions(tt.args)
			if got.WindowsMode != tt.want.WindowsMode {
				t.Errorf("WindowsMode = %v, want %v", got.WindowsMode, tt.want.WindowsMode)
			}
		})
	}
}

func TestWindowsFormat(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "ls-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files
	os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte("hello"), 0644)
	os.Mkdir(filepath.Join(tmpDir, "subdir"), 0755)

	// Test Windows format output
	output := FormatWindows(tmpDir, Options{})
	if !strings.Contains(output, "test.txt") {
		t.Error("output should contain test.txt")
	}
	if !strings.Contains(output, "subdir") {
		t.Error("output should contain subdir")
	}
	if !strings.Contains(output, "Directory of") {
		t.Error("output should contain 'Directory of'")
	}
}
