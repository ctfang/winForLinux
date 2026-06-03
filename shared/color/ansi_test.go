package color

import "testing"

func TestColorize(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		code     int
		expected string
	}{
		{"red", "hello", 31, "\033[31mhello\033[0m"},
		{"green", "file.exe", 32, "\033[32mfile.exe\033[0m"},
		{"blue", "dir", 34, "\033[34mdir\033[0m"},
		{"reset", "normal", 0, "\033[0mnormal\033[0m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Colorize(tt.text, tt.code)
			if result != tt.expected {
				t.Errorf("Colorize(%q, %d) = %q, want %q", tt.text, tt.code, result, tt.expected)
			}
		})
	}
}
