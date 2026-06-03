//go:build windows

package osutil

import (
	"os"
	"testing"
)

func TestIsExecutable(t *testing.T) {
	tests := []struct {
		name string
		file string
		want bool
	}{
		{"exe file", "test.exe", true},
		{"bat file", "test.bat", true},
		{"cmd file", "test.cmd", true},
		{"ps1 file", "test.ps1", true},
		{"com file", "test.com", true},
		{"msi file", "test.msi", true},
		{"txt file", "test.txt", false},
		{"go file", "test.go", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsExecutable(tt.file); got != tt.want {
				t.Errorf("IsExecutable(%q) = %v, want %v", tt.file, got, tt.want)
			}
		})
	}
}

func TestFormatSize(t *testing.T) {
	tests := []struct {
		name string
		size int64
		want string
	}{
		{"bytes", 500, "500B"},
		{"kilobytes", 1024, "1.0K"},
		{"megabytes", 1048576, "1.0M"},
		{"gigabytes", 1073741824, "1.0G"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatSize(tt.size); got != tt.want {
				t.Errorf("FormatSize(%d) = %q, want %q", tt.size, got, tt.want)
			}
		})
	}
}

func TestGetFileMode(t *testing.T) {
	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	mode := GetFileMode(tmpFile.Name())
	if len(mode) != 10 {
		t.Errorf("GetFileMode() length = %d, want 10", len(mode))
	}
	// Should be -rw-rw-rw- or similar format
	if mode[0] != '-' {
		t.Errorf("GetFileMode()[0] = %c, want '-'", mode[0])
	}
}
