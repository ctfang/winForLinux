package color

import (
	"os"
	"testing"
)

func TestIsTerminal(t *testing.T) {
	// In a test environment (go test), stdout is connected to a pipe, not a terminal.
	result := IsTerminal()
	if result {
		t.Errorf("IsTerminal() = true, want false (running under go test, stdout is a pipe)")
	}
}

func TestSupportsColor(t *testing.T) {
	// Save original env vars and restore after each subtest.
	origWT := os.Getenv("WT_SESSION")
	origConEmu := os.Getenv("ConEmuANSI")
	t.Cleanup(func() {
		os.Setenv("WT_SESSION", origWT)
		os.Setenv("ConEmuANSI", origConEmu)
	})

	t.Run("WT_Session_set", func(t *testing.T) {
		os.Setenv("WT_SESSION", "some-guid")
		os.Unsetenv("ConEmuANSI")
		if got := SupportsColor(); !got {
			t.Errorf("SupportsColor() = false with WT_SESSION set, want true")
		}
	})

	t.Run("ConEmuANSI_ON", func(t *testing.T) {
		os.Unsetenv("WT_SESSION")
		os.Setenv("ConEmuANSI", "ON")
		if got := SupportsColor(); !got {
			t.Errorf("SupportsColor() = false with ConEmuANSI=ON, want true")
		}
	})

	t.Run("neither_set_falls_back_to_IsTerminal", func(t *testing.T) {
		os.Unsetenv("WT_SESSION")
		os.Unsetenv("ConEmuANSI")
		want := IsTerminal()
		if got := SupportsColor(); got != want {
			t.Errorf("SupportsColor() = %v, want %v (should match IsTerminal())", got, want)
		}
	})
}
