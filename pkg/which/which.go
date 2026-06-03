package which

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Which finds an executable in PATH
func Which(name string) (string, error) {
	// If contains path separators, check directly
	if strings.ContainsAny(name, `/\`) {
		if _, err := os.Stat(name); err == nil {
			return name, nil
		}
		return "", fmt.Errorf("not found: %s", name)
	}

	// Get PATH environment variable
	pathEnv := os.Getenv("PATH")
	paths := strings.Split(pathEnv, ";")

	// Possible extensions
	extensions := []string{"", ".exe", ".bat", ".cmd", ".ps1"}

	for _, dir := range paths {
		for _, ext := range extensions {
			fullPath := filepath.Join(dir, name+ext)
			if _, err := os.Stat(fullPath); err == nil {
				return fullPath, nil
			}
		}
	}

	return "", fmt.Errorf("not found in PATH: %s", name)
}

// WhichAll finds all matching executables
func WhichAll(name string) ([]string, error) {
	var results []string

	if strings.ContainsAny(name, `/\`) {
		if _, err := os.Stat(name); err == nil {
			return []string{name}, nil
		}
		return nil, fmt.Errorf("not found: %s", name)
	}

	pathEnv := os.Getenv("PATH")
	paths := strings.Split(pathEnv, ";")
	extensions := []string{"", ".exe", ".bat", ".cmd", ".ps1"}

	for _, dir := range paths {
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
