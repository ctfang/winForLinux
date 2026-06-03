//go:build windows

package osutil

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"
)

var (
	kernel32                       = syscall.NewLazyDLL("kernel32.dll")
	procGetFileInformationByHandle = kernel32.NewProc("GetFileInformationByHandle")
)

// IsExecutable checks whether a file is an executable by extension.
func IsExecutable(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	switch ext {
	case ".exe", ".bat", ".cmd", ".ps1", ".com", ".msi":
		return true
	}
	return false
}

// FormatSize formats a file size into a human-readable string.
func FormatSize(size int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	switch {
	case size >= GB:
		return fmt.Sprintf("%.1fG", float64(size)/float64(GB))
	case size >= MB:
		return fmt.Sprintf("%.1fM", float64(size)/float64(MB))
	case size >= KB:
		return fmt.Sprintf("%.1fK", float64(size)/float64(KB))
	default:
		return fmt.Sprintf("%dB", size)
	}
}

// GetFileMode returns a Linux-style permission string (e.g. "-rw-rw-rw-").
func GetFileMode(name string) string {
	// Use Lstat so symlinks are detected instead of being followed.
	info, err := os.Lstat(name)
	if err != nil {
		return "?---------"
	}

	var b strings.Builder

	// File type character
	if info.Mode()&os.ModeSymlink != 0 {
		b.WriteRune('l')
	} else if info.IsDir() {
		b.WriteRune('d')
	} else {
		b.WriteRune('-')
	}

	// Build owner permissions (rwx)
	perm := info.Mode().Perm()
	if perm&0400 != 0 {
		b.WriteRune('r')
	} else {
		b.WriteRune('-')
	}
	if perm&0200 != 0 {
		b.WriteRune('w')
	} else {
		b.WriteRune('-')
	}
	if IsExecutable(name) {
		b.WriteRune('x')
	} else {
		b.WriteRune('-')
	}

	// On Windows, group and other permissions mirror owner.
	owner := b.String()[1:]
	b.WriteString(owner)
	b.WriteString(owner)

	return b.String()[:10]
}

// GetFileOwner returns the file owner (Windows: current user).
// The name parameter is accepted for API consistency but is not used on Windows.
func GetFileOwner(name string) string {
	return os.Getenv("USERNAME")
}

// GetFileGroup returns the file group (Windows: user domain).
// The name parameter is accepted for API consistency but is not used on Windows.
func GetFileGroup(name string) string {
	return os.Getenv("USERDOMAIN")
}

// IsHidden checks whether a file is hidden.
func IsHidden(name string) bool {
	base := filepath.Base(name)
	if strings.HasPrefix(base, ".") {
		return true
	}

	ptr, err := syscall.UTF16PtrFromString(name)
	if err != nil {
		return false
	}

	attrs, err := syscall.GetFileAttributes(ptr)
	if err != nil {
		return false
	}

	return attrs&syscall.FILE_ATTRIBUTE_HIDDEN != 0
}

// byHandleFileInformation mirrors the Windows BY_HANDLE_FILE_INFORMATION struct.
type byHandleFileInformation struct {
	FileAttributes     uint32
	CreationTime       syscall.Filetime
	LastAccessTime     syscall.Filetime
	LastWriteTime      syscall.Filetime
	VolumeSerialNumber uint32
	FileSizeHigh       uint32
	FileSizeLow        uint32
	NumberOfLinks      uint32
	FileIndexHigh      uint32
	FileIndexLow       uint32
}

// GetInode returns the file index on Windows (analogous to inode).
func GetInode(name string) uint64 {
	ptr, err := syscall.UTF16PtrFromString(name)
	if err != nil {
		return 0
	}
	handle, err := syscall.CreateFile(
		ptr,
		0,
		syscall.FILE_SHARE_READ|syscall.FILE_SHARE_WRITE|syscall.FILE_SHARE_DELETE,
		nil,
		syscall.OPEN_EXISTING,
		syscall.FILE_FLAG_BACKUP_SEMANTICS,
		0,
	)
	if err != nil {
		return 0
	}
	defer syscall.CloseHandle(handle)

	var info byHandleFileInformation

	ret, _, _ := procGetFileInformationByHandle.Call(
		uintptr(handle),
		uintptr(unsafe.Pointer(&info)),
	)

	if ret == 0 {
		return 0
	}

	return uint64(info.FileIndexHigh)<<32 | uint64(info.FileIndexLow)
}
