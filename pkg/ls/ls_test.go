package ls

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"
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

func TestLinuxFormat(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "ls-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte("hello"), 0644)
	os.Mkdir(filepath.Join(tmpDir, "subdir"), 0755)

	output, _ := FormatLinux(tmpDir, Options{Long: true})
	if !strings.Contains(output, "test.txt") {
		t.Error("output should contain test.txt")
	}
	if !strings.Contains(output, "subdir") {
		t.Error("output should contain subdir")
	}
}

func TestLinuxFormatAll(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "ls-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	os.WriteFile(filepath.Join(tmpDir, ".hidden"), []byte("hidden"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "visible"), []byte("visible"), 0644)

	// Without -a flag
	output, _ := FormatLinux(tmpDir, Options{})
	if strings.Contains(output, ".hidden") {
		t.Error("output should not contain .hidden without -a flag")
	}

	// With -a flag
	output, _ = FormatLinux(tmpDir, Options{All: true})
	if !strings.Contains(output, ".hidden") {
		t.Error("output should contain .hidden with -a flag")
	}
}

// TestLinuxLongFormatContent verifies that long format output contains the
// expected permission string, size, and modification time components.
func TestLinuxLongFormatContent(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "ls-test-long")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a file with known content length so size is predictable.
	data := strings.Repeat("x", 1024)
	err = os.WriteFile(filepath.Join(tmpDir, "bigfile.dat"), []byte(data), 0644)
	if err != nil {
		t.Fatal(err)
	}
	os.Mkdir(filepath.Join(tmpDir, "mydir"), 0755)

	output, err := FormatLinux(tmpDir, Options{Long: true})
	if err != nil {
		t.Fatal(err)
	}

	lines := strings.Split(strings.TrimRight(output, "\n"), "\n")
	if len(lines) < 2 {
		t.Fatalf("expected at least 2 lines, got %d: %q", len(lines), output)
	}

	// Permission string regex: starts with file type char, then 9 rwx/- chars.
	permRe := regexp.MustCompile(`^[dl\-][rwx\-]{9}`)

	foundFile := false
	foundDir := false
	for _, line := range lines {
		if strings.Contains(line, "bigfile.dat") {
			foundFile = true
			// Check permission string
			if !permRe.MatchString(line) {
				t.Errorf("long format line missing permission string: %q", line)
			}
			// Check that size 1024 appears
			if !strings.Contains(line, "1024") {
				t.Errorf("long format line should contain size 1024: %q", line)
			}
			// Check that a mod time pattern appears (e.g. "Jan  2 15:04")
			timeRe := regexp.MustCompile(`[A-Z][a-z]{2}\s+\d{1,2}\s+\d{2}:\d{2}`)
			if !timeRe.MatchString(line) {
				t.Errorf("long format line missing mod time: %q", line)
			}
		}
		if strings.Contains(line, "mydir") {
			foundDir = true
			if !permRe.MatchString(line) {
				t.Errorf("long format dir line missing permission string: %q", line)
			}
			// Directory permission string should start with 'd'
			if line[0] != 'd' {
				t.Errorf("directory permission should start with 'd', got %q", line)
			}
		}
	}

	if !foundFile {
		t.Error("output should contain bigfile.dat")
	}
	if !foundDir {
		t.Error("output should contain mydir")
	}
}

// TestLinuxLongFormatHumanSize verifies that -h flag produces human-readable sizes.
func TestLinuxLongFormatHumanSize(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "ls-test-human")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	data := strings.Repeat("x", 2048) // 2KB
	err = os.WriteFile(filepath.Join(tmpDir, "medium.dat"), []byte(data), 0644)
	if err != nil {
		t.Fatal(err)
	}

	output, err := FormatLinux(tmpDir, Options{Long: true, Human: true})
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(output, "medium.dat") {
		t.Error("output should contain medium.dat")
	}
	// Human-readable size should contain "K" for ~2KB
	if !strings.Contains(output, "K") {
		t.Errorf("human-readable output should contain 'K' for size: %q", output)
	}
}

// TestLinuxSortByTime verifies that -t flag sorts entries by modification time.
func TestLinuxSortByTime(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "ls-test-sort-time")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create files with staggered modification times.
	oldFile := filepath.Join(tmpDir, "old.txt")
	newFile := filepath.Join(tmpDir, "new.txt")
	os.WriteFile(oldFile, []byte("old"), 0644)
	// Ensure a measurable time difference.
	time.Sleep(50 * time.Millisecond)
	os.WriteFile(newFile, []byte("new"), 0644)

	output, err := FormatLinux(tmpDir, Options{TimeSort: true})
	if err != nil {
		t.Fatal(err)
	}

	// In time-sorted (newest first) output, "new.txt" should appear before "old.txt".
	idxNew := strings.Index(output, "new.txt")
	idxOld := strings.Index(output, "old.txt")
	if idxNew < 0 || idxOld < 0 {
		t.Fatalf("expected both files in output: %q", output)
	}
	if idxNew > idxOld {
		t.Errorf("time sort (newest first): new.txt (idx %d) should appear before old.txt (idx %d)", idxNew, idxOld)
	}

	// Reversed time sort: oldest first.
	output, err = FormatLinux(tmpDir, Options{TimeSort: true, Reverse: true})
	if err != nil {
		t.Fatal(err)
	}
	idxNew = strings.Index(output, "new.txt")
	idxOld = strings.Index(output, "old.txt")
	if idxNew < 0 || idxOld < 0 {
		t.Fatalf("expected both files in reversed output: %q", output)
	}
	if idxOld > idxNew {
		t.Errorf("reversed time sort (oldest first): old.txt (idx %d) should appear before new.txt (idx %d)", idxOld, idxNew)
	}
}

// TestLinuxSortBySize verifies that -S flag sorts entries by size (largest first).
func TestLinuxSortBySize(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "ls-test-sort-size")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	smallFile := filepath.Join(tmpDir, "small.txt")
	bigFile := filepath.Join(tmpDir, "big.txt")
	os.WriteFile(smallFile, []byte("hi"), 0644)
	os.WriteFile(bigFile, []byte(strings.Repeat("x", 4096)), 0644)

	output, err := FormatLinux(tmpDir, Options{SizeSort: true})
	if err != nil {
		t.Fatal(err)
	}

	idxBig := strings.Index(output, "big.txt")
	idxSmall := strings.Index(output, "small.txt")
	if idxBig < 0 || idxSmall < 0 {
		t.Fatalf("expected both files in output: %q", output)
	}
	if idxBig > idxSmall {
		t.Errorf("size sort (largest first): big.txt (idx %d) should appear before small.txt (idx %d)", idxBig, idxSmall)
	}

	// Reversed size sort: smallest first.
	output, err = FormatLinux(tmpDir, Options{SizeSort: true, Reverse: true})
	if err != nil {
		t.Fatal(err)
	}
	idxBig = strings.Index(output, "big.txt")
	idxSmall = strings.Index(output, "small.txt")
	if idxBig < 0 || idxSmall < 0 {
		t.Fatalf("expected both files in reversed output: %q", output)
	}
	if idxSmall > idxBig {
		t.Errorf("reversed size sort (smallest first): small.txt (idx %d) should appear before big.txt (idx %d)", idxSmall, idxBig)
	}
}

// TestLinuxOnePerLine verifies that -1 flag puts each entry on its own line.
func TestLinuxOnePerLine(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "ls-test-oneperline")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	os.WriteFile(filepath.Join(tmpDir, "a.txt"), []byte("a"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "b.txt"), []byte("b"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "c.txt"), []byte("c"), 0644)

	output, err := FormatLinux(tmpDir, Options{OnePerLine: true})
	if err != nil {
		t.Fatal(err)
	}

	lines := strings.Split(strings.TrimRight(output, "\n"), "\n")
	if len(lines) != 3 {
		t.Errorf("expected 3 lines for one-per-line output, got %d: %q", len(lines), output)
	}

	for _, name := range []string{"a.txt", "b.txt", "c.txt"} {
		found := false
		for _, line := range lines {
			if strings.TrimSpace(line) == name {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("one-per-line output missing %q in lines: %v", name, lines)
		}
	}
}

// TestLinuxSymlinkDetection verifies that symlinks are detected via Lstat
// and shown with 'l' permission prefix in long format.
func TestLinuxSymlinkDetection(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "ls-test-symlink")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a target file and a symlink to it.
	target := filepath.Join(tmpDir, "target.txt")
	os.WriteFile(target, []byte("target content"), 0644)
	link := filepath.Join(tmpDir, "link.txt")
	err = os.Symlink(target, link)
	if err != nil {
		t.Skipf("symlink creation not supported: %v", err)
	}

	// Long format should show 'l' prefix for symlinks.
	output, err := FormatLinux(tmpDir, Options{Long: true})
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(output, "link.txt") {
		t.Error("output should contain link.txt")
	}

	// Find the line for link.txt and check it starts with 'l'.
	for _, line := range strings.Split(output, "\n") {
		if strings.Contains(line, "link.txt") {
			if !strings.HasPrefix(line, "l") {
				t.Errorf("symlink line should start with 'l', got: %q", line)
			}
			break
		}
	}
}

// TestLinuxNonexistentDir verifies that listing a non-existent directory
// returns an error rather than panicking.
func TestLinuxNonexistentDir(t *testing.T) {
	_, err := FormatLinux("/nonexistent/path/that/does/not/exist", Options{})
	if err == nil {
		t.Error("expected error for non-existent directory")
	}
}

// TestSortEntriesErrorHandling verifies that sortEntries does not panic
// when entries have errors (e.g., deleted between readdir and stat).
func TestSortEntriesErrorHandling(t *testing.T) {
	// sortEntries with valid entries should not panic for any sort mode.
	tmpDir, err := os.MkdirTemp("", "ls-test-sort-err")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	os.WriteFile(filepath.Join(tmpDir, "a.txt"), []byte("aaa"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "b.txt"), []byte("b"), 0644)

	entries, err := os.ReadDir(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	// These should all complete without panicking.
	for _, opts := range []Options{
		{TimeSort: true},
		{SizeSort: true},
		{TimeSort: true, Reverse: true},
		{SizeSort: true, Reverse: true},
		{Reverse: true},
	} {
		sorted := sortEntries(entries, tmpDir, opts)
		if len(sorted) != len(entries) {
			t.Errorf("sortEntries returned %d entries, want %d", len(sorted), len(entries))
		}
	}
}
