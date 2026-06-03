package ls

import (
	"fmt"

	"github.com/user/win-tools/shared/argparse"
)

// Options ls command options
type Options struct {
	// Common options
	Paths []string

	// Windows mode options
	WindowsMode bool
	Wide        bool   // /W
	Bare        bool   // /B
	Recursive   bool   // /S
	SortBy      string // /O
	AttrFilter  string // /A

	// Linux mode options
	Long       bool   // -l
	All        bool   // -a, --all
	Human      bool   // -h
	TimeSort   bool   // -t
	Reverse    bool   // -r
	Color      string // --color
	OnePerLine bool   // -1
	DirOnly    bool   // -d
	FileType   bool   // -F
	SizeSort   bool   // -S
	Inode      bool   // -i
}

// ParseOptions parses options from arguments
func ParseOptions(args []string) Options {
	opts, positional := argparse.ParseArgs(args)

	options := Options{
		Paths: positional,
	}

	// Check if using Linux mode
	if argparse.HasLinuxOption(args) {
		// Linux mode
		options.Long = opts["l"] != "" || opts["long"] != ""
		options.All = opts["a"] != "" || opts["all"] != ""
		options.Human = opts["h"] != "" || opts["human-readable"] != ""
		options.TimeSort = opts["t"] != "" || opts["time"] != ""
		options.Reverse = opts["r"] != "" || opts["reverse"] != ""
		options.OnePerLine = opts["1"] != ""
		options.DirOnly = opts["d"] != ""
		options.FileType = opts["F"] != "" || opts["classify"] != ""
		options.SizeSort = opts["S"] != "" || opts["size"] != ""
		options.Inode = opts["i"] != "" || opts["inode"] != ""

		if v, ok := opts["color"]; ok {
			options.Color = v
		}
	} else {
		// Windows mode
		options.WindowsMode = true
		options.Wide = opts["W"] != ""
		options.Bare = opts["B"] != ""
		options.Recursive = opts["S"] != ""
		options.SortBy = opts["O"]
		options.AttrFilter = opts["A"]
	}

	// Default path is current directory
	if len(options.Paths) == 0 {
		options.Paths = []string{"."}
	}

	return options
}

// Run executes the ls command
func Run(args []string) error {
	options := ParseOptions(args)

	for _, path := range options.Paths {
		if options.WindowsMode {
			output := FormatWindows(path, options)
			fmt.Print(output)
		} else {
			// Linux mode will be implemented in Task 7
			output := FormatLinux(path, options)
			fmt.Print(output)
		}
	}

	return nil
}
