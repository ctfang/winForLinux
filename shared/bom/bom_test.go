package bom

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestReaderStripsBOM(t *testing.T) {
	in := append([]byte{0xEF, 0xBB, 0xBF}, []byte("alpha\nalpha\n")...)
	r := Reader(bytes.NewReader(in))
	b, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	got := string(b)
	if strings.HasPrefix(got, "\ufeff") {
		t.Fatalf("BOM not stripped: %q", got)
	}
	if got != "alpha\nalpha\n" {
		t.Fatalf("got %q", got)
	}
}

func TestReaderNoBOM(t *testing.T) {
	in := []byte("hello\n")
	r := Reader(bytes.NewReader(in))
	b, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "hello\n" {
		t.Fatalf("got %q", b)
	}
}
