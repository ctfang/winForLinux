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
	Recursive   bool   // /S or -R
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
		_, options.Long = opts["l"]
		if _, ok := opts["long"]; ok {
			options.Long = true
		}
		_, options.All = opts["a"]
		if _, ok := opts["all"]; ok {
			options.All = true
		}
		_, options.Human = opts["h"]
		if _, ok := opts["human-readable"]; ok {
			options.Human = true
		}
		_, options.TimeSort = opts["t"]
		if _, ok := opts["time"]; ok {
			options.TimeSort = true
		}
		_, options.Reverse = opts["r"]
		if _, ok := opts["reverse"]; ok {
			options.Reverse = true
		}
		_, options.Recursive = opts["R"]
		if _, ok := opts["recursive"]; ok {
			options.Recursive = true
		}
		_, options.OnePerLine = opts["1"]
		_, options.DirOnly = opts["d"]
		_, options.FileType = opts["F"]
		if _, ok := opts["classify"]; ok {
			options.FileType = true
		}
		_, options.SizeSort = opts["S"]
		if _, ok := opts["size"]; ok {
			options.SizeSort = true
		}
		_, options.Inode = opts["i"]
		if _, ok := opts["inode"]; ok {
			options.Inode = true
		}

		if v, ok := opts["color"]; ok {
			options.Color = v
		}
	} else {
		// Windows mode
		options.WindowsMode = true
		_, options.Wide = opts["W"]
		_, options.Bare = opts["B"]
		_, options.Recursive = opts["S"]
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
			output, err := FormatWindows(path, options)
			if err != nil {
				return err
			}
			fmt.Print(output)
		} else {
			output, err := FormatLinux(path, options)
			if err != nil {
				return err
			}
			fmt.Print(output)
		}
	}

	return nil
}
