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

func TestCatShowEnds(t *testing.T) {
	input := "line1\nline2\nline3"
	var buf bytes.Buffer
	err := Cat(strings.NewReader(input), Options{ShowEnds: true}, &buf)
	if err != nil {
		t.Fatalf("Cat() error = %v", err)
	}
	expected := "line1$\nline2$\nline3$\n"
	if buf.String() != expected {
		t.Errorf("Cat() = %q, want %q", buf.String(), expected)
	}
}

func TestCatShowTabs(t *testing.T) {
	input := "col1\tcol2\tcol3"
	var buf bytes.Buffer
	err := Cat(strings.NewReader(input), Options{ShowTabs: true}, &buf)
	if err != nil {
		t.Fatalf("Cat() error = %v", err)
	}
	expected := "col1^Icol2^Icol3\n"
	if buf.String() != expected {
		t.Errorf("Cat() = %q, want %q", buf.String(), expected)
	}
}

func TestCatSqueezeBlank(t *testing.T) {
	input := "line1\n\n\n\nline2\n\nline3"
	var buf bytes.Buffer
	err := Cat(strings.NewReader(input), Options{SqueezeBlank: true}, &buf)
	if err != nil {
		t.Fatalf("Cat() error = %v", err)
	}
	expected := "line1\n\nline2\n\nline3\n"
	if buf.String() != expected {
		t.Errorf("Cat() = %q, want %q", buf.String(), expected)
	}
}

func TestCatNumberSqueezeBlank(t *testing.T) {
	input := "line1\n\n\n\nline2\n\nline3"
	var buf bytes.Buffer
	err := Cat(strings.NewReader(input), Options{Number: true, SqueezeBlank: true}, &buf)
	if err != nil {
		t.Fatalf("Cat() error = %v", err)
	}
	expected := "     1\tline1\n     2\t\n     3\tline2\n     4\t\n     5\tline3\n"
	if buf.String() != expected {
		t.Errorf("Cat() = %q, want %q", buf.String(), expected)
	}
}
