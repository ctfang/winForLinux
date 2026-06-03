package argparse

import "strings"

// linuxOptions defines Linux-specific options that trigger Linux mode.
var linuxOptions = map[string]bool{
	"l": true, "a": true, "h": true, "t": true,
	"R": true, "r": true, "1": true, "d": true,
	"F": true, "S": true, "i": true,
	"all": true, "color": true, "human-readable": true,
}

// ParseArgs parses command-line arguments.
// Returns: (options map, positional arguments list)
func ParseArgs(args []string) (map[string]string, []string) {
	opts := make(map[string]string)
	var positional []string

	for i := 0; i < len(args); i++ {
		arg := args[i]

		if strings.HasPrefix(arg, "--") {
			// Long option: --option or --option=value
			key := arg[2:]
			if idx := strings.Index(key, "="); idx != -1 {
				opts[key[:idx]] = key[idx+1:]
			} else {
				opts[key] = ""
			}
		} else if strings.HasPrefix(arg, "-") && len(arg) > 1 {
			// Short option: -l, -a, -lah (combined)
			flags := arg[1:]
			for _, ch := range flags {
				opts[string(ch)] = ""
			}
		} else if strings.HasPrefix(arg, "/") && len(arg) >= 2 {
			// Windows style: /W, /B, /O:value, /A:value
			flag := arg[1:]
			if idx := strings.Index(flag, ":"); idx != -1 {
				opts[flag[:idx]] = flag[idx+1:]
			} else {
				opts[flag] = ""
			}
		} else {
			// Positional argument
			positional = append(positional, arg)
		}
	}

	return opts, positional
}

// HasLinuxOption checks whether args contain Linux-specific options.
func HasLinuxOption(args []string) bool {
	for _, arg := range args {
		if strings.HasPrefix(arg, "--") {
			key := arg[2:]
			if idx := strings.Index(key, "="); idx != -1 {
				key = key[:idx]
			}
			if linuxOptions[key] {
				return true
			}
		} else if strings.HasPrefix(arg, "-") && len(arg) > 1 {
			flags := arg[1:]
			for _, ch := range flags {
				if linuxOptions[string(ch)] {
					return true
				}
			}
		}
	}
	return false
}
