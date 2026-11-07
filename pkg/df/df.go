package df

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"syscall"
	"unsafe"
)

// Options for df
type Options struct {
	Human bool // -h
}

type usage struct {
	FS      string
	Total   uint64
	Free    uint64
	Avail   uint64
	Used    uint64
	UsePct  int
	Mounted string
}

// Df prints disk free space for paths (or all fixed drives).
func Df(paths []string, opts Options, w io.Writer) error {
	if len(paths) == 0 {
		paths = listDrives()
	}
	fmt.Fprintf(w, "%-10s %10s %10s %10s %6s %s\n", "Filesystem", sizeHeader(opts), "Used", "Avail", "Use%", "Mounted on")
	seen := map[string]bool{}
	for _, p := range paths {
		root := volumeRoot(p)
		if seen[root] {
			continue
		}
		seen[root] = true
		u, err := getUsage(root)
		if err != nil {
			fmt.Fprintf(os.Stderr, "df: %s: %v\n", p, err)
			continue
		}
		u.Mounted = root
		u.FS = root
		printUsage(u, opts, w)
	}
	return nil
}

func sizeHeader(opts Options) string {
	if opts.Human {
		return "Size"
	}
	return "1K-blocks"
}

func printUsage(u usage, opts Options, w io.Writer) {
	if opts.Human {
		fmt.Fprintf(w, "%-10s %10s %10s %10s %5d%% %s\n",
			u.FS, human(u.Total), human(u.Used), human(u.Avail), u.UsePct, u.Mounted)
		return
	}
	toK := func(b uint64) uint64 { return b / 1024 }
	fmt.Fprintf(w, "%-10s %10d %10d %10d %5d%% %s\n",
		u.FS, toK(u.Total), toK(u.Used), toK(u.Avail), u.UsePct, u.Mounted)
}

func human(n uint64) string {
	units := []string{"", "K", "M", "G", "T", "P"}
	v := float64(n)
	i := 0
	for v >= 1024 && i < len(units)-1 {
		v /= 1024
		i++
	}
	if i == 0 {
		return fmt.Sprintf("%d", n)
	}
	if v >= 10 {
		return fmt.Sprintf("%.0f%s", v, units[i])
	}
	return fmt.Sprintf("%.1f%s", v, units[i])
}

func volumeRoot(p string) string {
	abs, err := filepath.Abs(p)
	if err != nil {
		abs = p
	}
	vol := filepath.VolumeName(abs)
	if vol == "" {
		return abs
	}
	return vol + string(os.PathSeparator)
}

func listDrives() []string {
	var out []string
	for d := 'A'; d <= 'Z'; d++ {
		root := string(d) + ":\\"
		st, err := os.Stat(root)
		if err == nil && st.IsDir() {
			out = append(out, root)
		}
	}
	if len(out) == 0 {
		out = []string{"."}
	}
	return out
}

func getUsage(root string) (usage, error) {
	var u usage
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	proc := kernel32.NewProc("GetDiskFreeSpaceExW")
	pathPtr, err := syscall.UTF16PtrFromString(root)
	if err != nil {
		return u, err
	}
	var freeBytesAvailable, totalNumberOfBytes, totalNumberOfFreeBytes uint64
	r1, _, callErr := proc.Call(
		uintptr(unsafe.Pointer(pathPtr)),
		uintptr(unsafe.Pointer(&freeBytesAvailable)),
		uintptr(unsafe.Pointer(&totalNumberOfBytes)),
		uintptr(unsafe.Pointer(&totalNumberOfFreeBytes)),
	)
	if r1 == 0 {
		return u, fmt.Errorf("%v", callErr)
	}
	u.Total = totalNumberOfBytes
	u.Avail = freeBytesAvailable
	u.Free = totalNumberOfFreeBytes
	if u.Total >= u.Avail {
		u.Used = u.Total - u.Avail
	}
	if u.Total > 0 {
		u.UsePct = int((u.Used * 100) / u.Total)
	}
	return u, nil
}

