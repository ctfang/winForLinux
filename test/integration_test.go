package test

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func getBinaryPath(name string) string {
	if runtime.GOOS == "windows" {
		return filepath.Join("..", "dist", name+".exe")
	}
	return filepath.Join("..", "dist", name)
}

func TestLsIntegration(t *testing.T) {
	lsPath := getBinaryPath("ls")
	if _, err := os.Stat(lsPath); os.IsNotExist(err) {
		t.Skip("ls binary not found, run build.ps1 first")
	}

	// Test Windows mode
	out, err := exec.Command(lsPath).CombinedOutput()
	if err != nil {
		t.Fatalf("ls failed: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "Directory of") {
		t.Error("Windows mode should contain 'Directory of'")
	}

	// Test Linux mode
	out, err = exec.Command(lsPath, "-l").CombinedOutput()
	if err != nil {
		t.Fatalf("ls -l failed: %v\n%s", err, out)
	}
}

func TestGrepIntegration(t *testing.T) {
	grepPath := getBinaryPath("grep")
	if _, err := os.Stat(grepPath); os.IsNotExist(err) {
		t.Skip("grep binary not found, run build.ps1 first")
	}

	// Test stdin
	cmd := exec.Command(grepPath, "hello")
	cmd.Stdin = strings.NewReader("hello\nworld\nhello world")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("grep failed: %v\n%s", err, out)
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 matches, got %d", len(lines))
	}
}

func TestCatIntegration(t *testing.T) {
	catPath := getBinaryPath("cat")
	if _, err := os.Stat(catPath); os.IsNotExist(err) {
		t.Skip("cat binary not found, run build.ps1 first")
	}

	// Create temp file
	tmpFile, err := os.CreateTemp("", "cat-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.WriteString("line1\nline2\n")
	tmpFile.Close()

	out, err := exec.Command(catPath, "-n", tmpFile.Name()).CombinedOutput()
	if err != nil {
		t.Fatalf("cat -n failed: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "1") || !strings.Contains(string(out), "line1") {
		t.Error("cat -n should show line numbers")
	}
}

func TestWhichIntegration(t *testing.T) {
	whichPath := getBinaryPath("which")
	if _, err := os.Stat(whichPath); os.IsNotExist(err) {
		t.Skip("which binary not found, run build.ps1 first")
	}

	out, err := exec.Command(whichPath, "go").CombinedOutput()
	if err != nil {
		t.Fatalf("which go failed: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "go") {
		t.Error("which go should return go path")
	}
}
