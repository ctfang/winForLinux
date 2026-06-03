package which

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestWhich(t *testing.T) {
	// Test finding a known existing command
	result, err := Which("go")
	if err != nil {
		t.Fatalf("Which() error = %v", err)
	}
	if result == "" {
		t.Error("Which('go') should return a path")
	}
	t.Logf("go found at: %s", result)
}

func TestWhichNotFound(t *testing.T) {
	_, err := Which("nonexistent_command_12345")
	if err == nil {
		t.Error("Which() should return error for nonexistent command")
	}
}

func TestWhichAll(t *testing.T) {
	// Test finding a known existing command
	results, err := WhichAll("go")
	if err != nil {
		t.Fatalf("WhichAll() error = %v", err)
	}
	if len(results) == 0 {
		t.Error("WhichAll('go') should return at least one path")
	}
	t.Logf("go found at: %v", results)
}

func TestWhichAllNotFound(t *testing.T) {
	_, err := WhichAll("nonexistent_command_12345")
	if err == nil {
		t.Error("WhichAll() should return error for nonexistent command")
	}
}

func TestWhichAllMultipleResults(t *testing.T) {
	// Create a temporary directory structure with a duplicate executable
	tmpDir := t.TempDir()

	// Create a fake executable
	binName := "faketestcmd"
	if runtime.GOOS == "windows" {
		binName += ".exe"
	}
	binPath := filepath.Join(tmpDir, binName)
	if err := os.WriteFile(binPath, []byte("#!/bin/sh\n"), 0755); err != nil {
		t.Fatalf("failed to create fake binary: %v", err)
	}

	// Temporarily override PATH to include our temp dir twice
	origPath := os.Getenv("PATH")
	defer os.Setenv("PATH", origPath)

	dualDir := tmpDir + string(os.PathListSeparator) + tmpDir
	os.Setenv("PATH", dualDir)

	results, err := WhichAll("faketestcmd")
	if err != nil {
		t.Fatalf("WhichAll() error = %v", err)
	}
	// Both entries in PATH point to the same directory, so we should find 2 results
	if len(results) != 2 {
		t.Errorf("WhichAll() returned %d results, want 2: %v", len(results), results)
	}
}

func TestWhichAllEmptyPathSegments(t *testing.T) {
	// Test that empty PATH segments are filtered out (they don't cause errors)
	origPath := os.Getenv("PATH")
	defer os.Setenv("PATH", origPath)

	// Create a temp dir with a fake executable
	tmpDir := t.TempDir()
	binName := "emptysegcmd"
	if runtime.GOOS == "windows" {
		binName += ".exe"
	}
	binPath := filepath.Join(tmpDir, binName)
	if err := os.WriteFile(binPath, []byte("#!/bin/sh\n"), 0755); err != nil {
		t.Fatalf("failed to create fake binary: %v", err)
	}

	// Set PATH with empty segments (leading, trailing, and middle) around the real dir
	sep := string(os.PathListSeparator)
	os.Setenv("PATH", sep+tmpDir+sep+sep+tmpDir+sep)

	results, err := WhichAll("emptysegcmd")
	if err != nil {
		t.Fatalf("WhichAll() error = %v", err)
	}
	// The directory appears twice, so we expect 2 results (not broken by empty segments)
	if len(results) != 2 {
		t.Errorf("WhichAll() returned %d results, want 2: %v", len(results), results)
	}
}

func TestWhichAllDirectPath(t *testing.T) {
	// Test finding a command by direct path
	tmpDir := t.TempDir()
	binName := "directcmd"
	if runtime.GOOS == "windows" {
		binName += ".exe"
	}
	binPath := filepath.Join(tmpDir, binName)
	if err := os.WriteFile(binPath, []byte("#!/bin/sh\n"), 0755); err != nil {
		t.Fatalf("failed to create fake binary: %v", err)
	}

	results, err := WhichAll(binPath)
	if err != nil {
		t.Fatalf("WhichAll() error = %v", err)
	}
	if len(results) != 1 {
		t.Errorf("WhichAll() returned %d results, want 1: %v", len(results), results)
	}
	if results[0] != binPath {
		t.Errorf("WhichAll() = %s, want %s", results[0], binPath)
	}
}

func TestWhichAllDirectPathNotFound(t *testing.T) {
	_, err := WhichAll("/nonexistent/path/to/cmd")
	if err == nil {
		t.Error("WhichAll() should return error for nonexistent direct path")
	}
}

func TestWhichDelegatesToWhichAll(t *testing.T) {
	// Test that Which returns the first result from WhichAll
	whichResult, whichErr := Which("go")
	allResults, allErr := WhichAll("go")

	if whichErr != nil {
		t.Fatalf("Which() error = %v", whichErr)
	}
	if allErr != nil {
		t.Fatalf("WhichAll() error = %v", allErr)
	}
	if whichResult != allResults[0] {
		t.Errorf("Which() = %s, WhichAll()[0] = %s; should be equal", whichResult, allResults[0])
	}
}

func TestGetExtensions(t *testing.T) {
	exts := getExtensions()
	if len(exts) == 0 {
		t.Error("getExtensions() should return at least one element")
	}
	// First element should always be empty string (for bare names)
	if exts[0] != "" {
		t.Errorf("getExtensions()[0] = %q, want empty string", exts[0])
	}

	if runtime.GOOS == "windows" {
		// On Windows, should have extensions from PATHEXT
		if len(exts) < 2 {
			t.Error("getExtensions() on Windows should return extensions from PATHEXT")
		}
	} else {
		// On non-Windows, should only have the empty string
		if len(exts) != 1 {
			t.Errorf("getExtensions() on non-Windows returned %d elements, want 1", len(exts))
		}
	}
}

func TestGetPathDirs(t *testing.T) {
	origPath := os.Getenv("PATH")
	defer os.Setenv("PATH", origPath)

	// Test with empty segments
	os.Setenv("PATH", "a"+string(os.PathListSeparator)+""+string(os.PathListSeparator)+"b")
	dirs := getPathDirs()
	if len(dirs) != 2 {
		t.Errorf("getPathDirs() returned %d dirs, want 2: %v", len(dirs), dirs)
	}
}
