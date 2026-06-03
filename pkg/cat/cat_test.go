package cat

import (
	"bytes"
	"strings"
	"testing"
)

func TestCat(t *testing.T) {
	input := "line1\nline2\nline3\n"
	var buf bytes.Buffer
	err := Cat(strings.NewReader(input), Options{}, &buf)
	if err != nil {
		t.Fatalf("Cat() error = %v", err)
	}
	if buf.String() != input {
		t.Errorf("Cat() = %q, want %q", buf.String(), input)
	}
}

func TestCatNumber(t *testing.T) {
	input := "line1\nline2\nline3"
	var buf bytes.Buffer
	err := Cat(strings.NewReader(input), Options{Number: true}, &buf)
	if err != nil {
		t.Fatalf("Cat() error = %v", err)
	}
	expected := "     1\tline1\n     2\tline2\n     3\tline3\n"
	if buf.String() != expected {
		t.Errorf("Cat() = %q, want %q", buf.String(), expected)
	}
}
