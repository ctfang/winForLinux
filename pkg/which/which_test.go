package which

import "testing"

func TestWhich(t *testing.T) {
	// Test finding a known existing command
	result, err := Which("go")
	if err != nil {
		t.Fatalf("Which() error = %v", err)
	}
	if result == "" {
		t.Error("Which('go') should return a path")
	}
	t.Logf("go found at: %s", result)
}

func TestWhichNotFound(t *testing.T) {
	_, err := Which("nonexistent_command_12345")
	if err == nil {
		t.Error("Which() should return error for nonexistent command")
	}
}
