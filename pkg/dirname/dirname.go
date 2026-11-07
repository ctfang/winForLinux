package dirname

import "path/filepath"

// Dir returns the directory portion of path.
func Dir(path string) string {
	d := filepath.Dir(path)
	if d == "" {
		return "."
	}
	return d
}
