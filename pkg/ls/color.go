package ls

import (
	"path/filepath"
	"strings"

	"github.com/user/win-tools/shared/color"
	"github.com/user/win-tools/shared/osutil"
)

// typeIndicator returns the file type suffix character for -F (--classify).
// Returns "/" for directories, "*" for executables, "@" for symlinks, and
// "" for other file types.
func typeIndicator(name string, isDir bool, isSymlink bool) string {
	if isDir {
		return "/"
	} else if osutil.IsExecutable(name) {
		return "*"
	} else if isSymlink {
		return "@"
	}
	return ""
}

// colorizeName applies color to filename based on type
func colorizeName(name string, isDir bool, isSymlink bool, opts Options) string {
	if opts.Color == "never" || (opts.Color == "auto" && !color.SupportsColor()) {
		if opts.FileType {
			return name + typeIndicator(name, isDir, isSymlink)
		}
		return name
	}

	var result string
	if isDir {
		result = color.Colorize(name, color.Blue)
	} else if osutil.IsExecutable(name) {
		result = color.Colorize(name, color.Green)
	} else if isSymlink {
		result = color.Colorize(name, color.Cyan)
	} else {
		ext := strings.ToLower(filepath.Ext(name))
		switch ext {
		case ".zip", ".gz", ".tar", ".rar", ".7z", ".bz2":
			result = color.Colorize(name, color.Red)
		case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".svg":
			result = color.Colorize(name, color.Magenta)
		default:
			result = name
		}
	}

	if opts.FileType {
		result += typeIndicator(name, isDir, isSymlink)
	}

	return result
}
