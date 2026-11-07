package basename

import (
	"path/filepath"
	"strings"
)

// Base returns the base name of path, optionally stripping suffix.
func Base(path, suffix string) string {
	base := filepath.Base(path)
	if suffix != "" && strings.HasSuffix(base, suffix) && base != suffix {
		base = strings.TrimSuffix(base, suffix)
	}
	return base
}
