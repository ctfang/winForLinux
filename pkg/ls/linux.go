package ls

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/user/win-tools/shared/osutil"
)

// FormatLinux generates Linux ls style output
func FormatLinux(dir string, opts Options) (string, error) {
	var sb strings.Builder

	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", fmt.Errorf("ls: cannot access '%s': %w", dir, err)
	}

	// Filter hidden files
	var filtered []os.DirEntry
	for _, entry := range entries {
		if !opts.All && strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		filtered = append(filtered, entry)
	}

	// Sort entries
	sortedEntries := sortEntries(filtered, dir, opts)

	// Output
	for _, entry := range sortedEntries {
		fullPath := filepath.Join(dir, entry.Name())
		// Use Lstat instead of entry.Info() so symlinks are not followed.
		// entry.Info() calls os.Stat which follows symlinks, so ModeSymlink
		// would never be set.
		info, err := os.Lstat(fullPath)
		if err != nil {
			continue
		}

		if opts.Inode {
			sb.WriteString(fmt.Sprintf("%d ", osutil.GetInode(fullPath)))
		}

		isSymlink := info.Mode()&os.ModeSymlink != 0

		if opts.Long {
			sb.WriteString(formatLinuxLong(entry.Name(), info, fullPath, opts))
		} else if opts.OnePerLine {
			sb.WriteString(colorizeName(entry.Name(), info.IsDir(), isSymlink, opts) + "\n")
		} else {
			sb.WriteString(colorizeName(entry.Name(), info.IsDir(), isSymlink, opts) + "  ")
		}
	}

	if !opts.Long && !opts.OnePerLine {
		sb.WriteString("\n")
	}

	return sb.String(), nil
}

func sortEntries(entries []os.DirEntry, dir string, opts Options) []os.DirEntry {
	sorted := make([]os.DirEntry, len(entries))
	copy(sorted, entries)

	sort.Slice(sorted, func(i, j int) bool {
		if opts.TimeSort {
			infoI, errI := sorted[i].Info()
			infoJ, errJ := sorted[j].Info()
			// If either entry has an error, fall back to name comparison
			// to avoid nil pointer panics.
			if errI != nil || errJ != nil {
				if opts.Reverse {
					return sorted[i].Name() > sorted[j].Name()
				}
				return sorted[i].Name() < sorted[j].Name()
			}
			if opts.Reverse {
				return infoI.ModTime().Before(infoJ.ModTime())
			}
			return infoI.ModTime().After(infoJ.ModTime())
		}
		if opts.SizeSort {
			infoI, errI := sorted[i].Info()
			infoJ, errJ := sorted[j].Info()
			if errI != nil || errJ != nil {
				if opts.Reverse {
					return sorted[i].Name() > sorted[j].Name()
				}
				return sorted[i].Name() < sorted[j].Name()
			}
			if opts.Reverse {
				return infoI.Size() < infoJ.Size()
			}
			return infoI.Size() > infoJ.Size()
		}
		// Default: sort by name
		if opts.Reverse {
			return sorted[i].Name() > sorted[j].Name()
		}
		return sorted[i].Name() < sorted[j].Name()
	})

	return sorted
}

func formatLinuxLong(name string, info os.FileInfo, fullPath string, opts Options) string {
	var sb strings.Builder

	// Permissions
	mode := osutil.GetFileMode(fullPath)
	sb.WriteString(mode + " ")

	// Links
	sb.WriteString("1 ")

	// Owner
	sb.WriteString(osutil.GetFileOwner(fullPath) + " ")

	// Group
	sb.WriteString(osutil.GetFileGroup(fullPath) + " ")

	// Size
	if opts.Human {
		sb.WriteString(fmt.Sprintf("%6s ", osutil.FormatSize(info.Size())))
	} else {
		sb.WriteString(fmt.Sprintf("%6d ", info.Size()))
	}

	// Mod time
	sb.WriteString(info.ModTime().Format("Jan  2 15:04") + " ")

	// Name (with color)
	sb.WriteString(colorizeName(name, info.IsDir(), info.Mode()&os.ModeSymlink != 0, opts))

	sb.WriteString("\n")

	return sb.String()
}
