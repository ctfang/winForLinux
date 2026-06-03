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
			want: Options{WindowsMode: true, Wide: true, Bare: true, Paths: []string{"."}},
		},
		{
			name: "windows mode with values",
			args: []string{"/O:gn", "/A:d"},
			want: Options{WindowsMode: true, SortBy: "gn", AttrFilter: "d", Paths: []string{"."}},
		},
		{
			name: "linux mode",
			args: []string{"-l", "-a"},
			want: Options{Long: true, All: true, Paths: []string{"."}},
		},
		{
			name: "linux mode human-readable",
			args: []string{"-l", "-h"},
			want: Options{Long: true, Human: true, Paths: []string{"."}},
		},
		{
			name: "linux mode --human-readable",
			args: []string{"--human-readable"},
			want: Options{Human: true, Paths: []string{"."}},
		},
		{
			name: "no args",
			args: []string{},
			want: Options{WindowsMode: true, Paths: []string{"."}},
		},
		{
			name: "windows recursive with paths",
			args: []string{"/S", "/W", "dir1", "dir2"},
			want: Options{WindowsMode: true, Recursive: true, Wide: true, Paths: []string{"dir1", "dir2"}},
		},
		{
			name: "linux mode with path",
			args: []string{"-l", "mydir"},
			want: Options{Long: true, Paths: []string{"mydir"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseOptions(tt.args)
			if got.WindowsMode != tt.want.WindowsMode {
				t.Errorf("WindowsMode = %v, want %v", got.WindowsMode, tt.want.WindowsMode)
			}
			if got.Wide != tt.want.Wide {
				t.Errorf("Wide = %v, want %v", got.Wide, tt.want.Wide)
			}
			if got.Bare != tt.want.Bare {
				t.Errorf("Bare = %v, want %v", got.Bare, tt.want.Bare)
			}
			if got.Recursive != tt.want.Recursive {
				t.Errorf("Recursive = %v, want %v", got.Recursive, tt.want.Recursive)
			}
			if got.SortBy != tt.want.SortBy {
				t.Errorf("SortBy = %q, want %q", got.SortBy, tt.want.SortBy)
			}
			if got.AttrFilter != tt.want.AttrFilter {
				t.Errorf("AttrFilter = %q, want %q", got.AttrFilter, tt.want.AttrFilter)
			}
			if got.Long != tt.want.Long {
				t.Errorf("Long = %v, want %v", got.Long, tt.want.Long)
			}
			if got.All != tt.want.All {
				t.Errorf("All = %v, want %v", got.All, tt.want.All)
			}
			if got.Human != tt.want.Human {
				t.Errorf("Human = %v, want %v", got.Human, tt.want.Human)
			}
			if got.TimeSort != tt.want.TimeSort {
				t.Errorf("TimeSort = %v, want %v", got.TimeSort, tt.want.TimeSort)
			}
			if got.Reverse != tt.want.Reverse {
				t.Errorf("Reverse = %v, want %v", got.Reverse, tt.want.Reverse)
			}
			if got.OnePerLine != tt.want.OnePerLine {
				t.Errorf("OnePerLine = %v, want %v", got.OnePerLine, tt.want.OnePerLine)
			}
			if got.DirOnly != tt.want.DirOnly {
				t.Errorf("DirOnly = %v, want %v", got.DirOnly, tt.want.DirOnly)
			}
			if got.FileType != tt.want.FileType {
				t.Errorf("FileType = %v, want %v", got.FileType, tt.want.FileType)
			}
			if got.SizeSort != tt.want.SizeSort {
				t.Errorf("SizeSort = %v, want %v", got.SizeSort, tt.want.SizeSort)
			}
			if got.Inode != tt.want.Inode {
				t.Errorf("Inode = %v, want %v", got.Inode, tt.want.Inode)
			}
			if len(got.Paths) != len(tt.want.Paths) {
				t.Errorf("Paths = %v, want %v", got.Paths, tt.want.Paths)
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
	output, err := FormatWindows(tmpDir, Options{})
	if err != nil {
		t.Fatal(err)
	}
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
