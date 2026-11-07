package winproc

import (
	"fmt"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"
)

const (
	th32csSnapProcess   = 0x00000002
	processTerminate    = 0x0001
	processQueryLimited = 0x1000
)

var (
	kernel32                 = syscall.NewLazyDLL("kernel32.dll")
	procCreateToolhelp32Snap = kernel32.NewProc("CreateToolhelp32Snapshot")
	procProcess32FirstW      = kernel32.NewProc("Process32FirstW")
	procProcess32NextW       = kernel32.NewProc("Process32NextW")
	procOpenProcess          = kernel32.NewProc("OpenProcess")
	procTerminateProcess     = kernel32.NewProc("TerminateProcess")
	procCloseHandle          = kernel32.NewProc("CloseHandle")
	procQueryFullProcessName = kernel32.NewProc("QueryFullProcessImageNameW")
)

type processEntry32W struct {
	Size            uint32
	CntUsage        uint32
	ProcessID       uint32
	DefaultHeapID   uintptr
	ModuleID        uint32
	CntThreads      uint32
	ParentProcessID uint32
	PriClassBase    int32
	Flags           uint32
	ExeFile         [260]uint16
}

// Process is a running process snapshot entry.
type Process struct {
	PID     uint32
	PPID    uint32
	Threads uint32
	Name    string
	ExePath string
}

// List returns all processes visible to the caller.
func List() ([]Process, error) {
	snap, _, err := procCreateToolhelp32Snap.Call(th32csSnapProcess, 0)
	if snap == uintptr(syscall.InvalidHandle) {
		return nil, fmt.Errorf("CreateToolhelp32Snapshot: %v", err)
	}
	defer procCloseHandle.Call(snap)

	var pe processEntry32W
	pe.Size = uint32(unsafe.Sizeof(pe))
	r1, _, err := procProcess32FirstW.Call(snap, uintptr(unsafe.Pointer(&pe)))
	if r1 == 0 {
		return nil, fmt.Errorf("Process32First: %v", err)
	}

	var out []Process
	for {
		name := syscall.UTF16ToString(pe.ExeFile[:])
		p := Process{
			PID:     pe.ProcessID,
			PPID:    pe.ParentProcessID,
			Threads: pe.CntThreads,
			Name:    name,
		}
		p.ExePath = imagePath(pe.ProcessID)
		out = append(out, p)

		pe.Size = uint32(unsafe.Sizeof(pe))
		r1, _, _ = procProcess32NextW.Call(snap, uintptr(unsafe.Pointer(&pe)))
		if r1 == 0 {
			break
		}
	}
	return out, nil
}

// ByPID returns a single process by PID, or nil if not found.
func ByPID(pid uint32) (*Process, error) {
	all, err := List()
	if err != nil {
		return nil, err
	}
	for i := range all {
		if all[i].PID == pid {
			return &all[i], nil
		}
	}
	return nil, nil
}

// ByName returns processes whose exe name matches (case-insensitive).
func ByName(name string) ([]Process, error) {
	want := strings.ToLower(name)
	wantBase := strings.ToLower(strings.TrimSuffix(name, ".exe"))
	all, err := List()
	if err != nil {
		return nil, err
	}
	var out []Process
	for _, p := range all {
		n := strings.ToLower(p.Name)
		base := strings.ToLower(strings.TrimSuffix(p.Name, ".exe"))
		if n == want || base == wantBase || n == want+".exe" {
			out = append(out, p)
		}
	}
	return out, nil
}

// NameOf returns the executable base name for pid, or empty string.
func NameOf(pid uint32) string {
	if pid == 0 {
		return ""
	}
	p, err := ByPID(pid)
	if err != nil || p == nil {
		return ""
	}
	return p.Name
}

// Kill terminates a process. force is accepted for API compatibility (-9).
func Kill(pid uint32, force bool) error {
	if pid == 0 {
		return fmt.Errorf("invalid pid 0")
	}
	h, _, err := procOpenProcess.Call(uintptr(processTerminate), 0, uintptr(pid))
	if h == 0 {
		return fmt.Errorf("OpenProcess(%d): %v", pid, err)
	}
	defer procCloseHandle.Call(h)
	r1, _, err := procTerminateProcess.Call(h, 1)
	if r1 == 0 {
		return fmt.Errorf("TerminateProcess(%d): %v", pid, err)
	}
	_ = force
	return nil
}

func imagePath(pid uint32) string {
	h, _, _ := procOpenProcess.Call(uintptr(processQueryLimited), 0, uintptr(pid))
	if h == 0 {
		return ""
	}
	defer procCloseHandle.Call(h)
	buf := make([]uint16, 1024)
	size := uint32(len(buf))
	r1, _, _ := procQueryFullProcessName.Call(h, 0, uintptr(unsafe.Pointer(&buf[0])), uintptr(unsafe.Pointer(&size)))
	if r1 == 0 {
		return ""
	}
	return syscall.UTF16ToString(buf[:size])
}

// DisplayName returns basename of ExePath or Name.
func (p Process) DisplayName() string {
	if p.ExePath != "" {
		return filepath.Base(p.ExePath)
	}
	return p.Name
}