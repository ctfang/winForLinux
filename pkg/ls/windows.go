package ls

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// FormatWindows generates Windows dir style output
func FormatWindows(dir string, opts Options) string {
	var sb strings.Builder

	absDir, err := filepath.Abs(dir)
	if err != nil {
		absDir = dir
	}

	// Get drive info
	drive := filepath.VolumeName(absDir)
	if drive == "" {
		drive = "C:"
	}

	sb.WriteString(fmt.Sprintf(" Volume in drive %s has no label\n", string(drive[0])))
	sb.WriteString(" Volume Serial Number is 0000-0000\n")
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf(" Directory of %s\n\n", absDir))

	// Read directory
	entries, err := os.ReadDir(dir)
	if err != nil {
		sb.WriteString(fmt.Sprintf("Error reading directory: %v\n", err))
		return sb.String()
	}

	fileCount := 0
	dirCount := 0
	totalSize := int64(0)

	// Add . and ..
	sb.WriteString(formatWindowsEntry(".", true, 0, time.Time{}))
	sb.WriteString(formatWindowsEntry("..", true, 0, time.Time{}))
	dirCount += 2

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		isDir := entry.IsDir()
		if isDir {
			dirCount++
		} else {
			fileCount++
			totalSize += info.Size()
		}

		if opts.Bare {
			sb.WriteString(entry.Name() + "\n")
		} else if opts.Wide {
			if isDir {
				sb.WriteString(fmt.Sprintf("[%s]", entry.Name()))
			} else {
				sb.WriteString(entry.Name())
			}
			sb.WriteString("\t")
		} else {
			sb.WriteString(formatWindowsEntry(entry.Name(), isDir, info.Size(), info.ModTime()))
		}
	}

	if opts.Wide {
		sb.WriteString("\n")
	}

	sb.WriteString(fmt.Sprintf("%15d File(s) %15d bytes\n", fileCount, totalSize))
	sb.WriteString(fmt.Sprintf("%15d Dir(s)  %15d bytes free\n", dirCount, 0))

	return sb.String()
}

func formatWindowsEntry(name string, isDir bool, size int64, modTime time.Time) string {
	if isDir {
		if modTime.IsZero() {
			return fmt.Sprintf("               <DIR>          %s\n", name)
		}
		return fmt.Sprintf("%s  %s    <DIR>          %s\n",
			modTime.Format("2006/01/02"), modTime.Format("15:04"), name)
	}
	return fmt.Sprintf("%s  %s %15d %s\n",
		modTime.Format("2006/01/02"), modTime.Format("15:04"), size, name)
}
