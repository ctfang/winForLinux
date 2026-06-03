// Package which provides functions to locate executables in PATH,
// similar to the Unix which command.
package which

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// getExtensions returns the list of file extensions to try when searching
// for executables. On Windows, it reads the PATHEXT environment variable;
// on other platforms, it returns an empty list (executables have no extension).
func getExtensions() []string {
	if runtime.GOOS == "windows" {
		pathext := os.Getenv("PATHEXT")
		if pathext == "" {
			// Default Windows extensions if PATHEXT is not set
			return []string{".exe", ".bat", ".cmd", ".ps1"}
		}
		exts := strings.Split(pathext, string(os.PathListSeparator))
		// Always include empty string so bare names match first
		result := make([]string, 0, len(exts)+1)
		result = append(result, "")
		for _, ext := range exts {
			ext = strings.TrimSpace(ext)
			if ext != "" {
				result = append(result, ext)
			}
		}
		return result
	}
	// On non-Windows platforms, executables don't need an extension
	return []string{""}
}

// getPathDirs returns the list of directories from the PATH environment
// variable, filtering out empty segments.
func getPathDirs() []string {
	pathEnv := os.Getenv("PATH")
	parts := strings.Split(pathEnv, string(os.PathListSeparator))
	var dirs []string
	for _, p := range parts {
		if p != "" {
			dirs = append(dirs, p)
		}
	}
	return dirs
}

// WhichAll finds all matching executables for name in PATH.
// If name contains path separators, it is checked directly as a path.
func WhichAll(name string) ([]string, error) {
	// If contains path separators, check directly
	if strings.ContainsAny(name, `/\`) {
		if _, err := os.Stat(name); err == nil {
			return []string{name}, nil
		}
		return nil, fmt.Errorf("not found: %s", name)
	}

	dirs := getPathDirs()
	extensions := getExtensions()

	var results []string
	for _, dir := range dirs {
		for _, ext := range extensions {
			fullPath := filepath.Join(dir, name+ext)
			if _, err := os.Stat(fullPath); err == nil {
				results = append(results, fullPath)
			}
		}
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("not found in PATH: %s", name)
	}

	return results, nil
}

// Which finds the first matching executable for name in PATH.
// It is equivalent to calling WhichAll and returning the first result.
func Which(name string) (string, error) {
	results, err := WhichAll(name)
	if err != nil {
		return "", err
	}
	return results[0], nil
}
