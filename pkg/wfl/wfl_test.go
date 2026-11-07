package wfl

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
)

func TestTokenize(t *testing.T) {
	toks, err := Tokenize(`cat a.txt | grep -i "hello world" && echo ok`)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"cat", "a.txt", "|", "grep", "-i", "hello world", "&&", "echo", "ok"}
	if !reflect.DeepEqual(toks, want) {
		t.Fatalf("got %#v want %#v", toks, want)
	}
}

func TestParse(t *testing.T) {
	toks := []string{"false", "||", "echo", "fallback"}
	specs, err := Parse(toks)
	if err != nil {
		t.Fatal(err)
	}
	if len(specs) != 2 {
		t.Fatalf("len=%d", len(specs))
	}
	if specs[0].Op != OpOr || specs[1].Op != OpNone {
		t.Fatalf("ops: %v %v", specs[0].Op, specs[1].Op)
	}
}

func TestRunEchoPipe(t *testing.T) {
	var out bytes.Buffer
	code := RunLine("echo hello-world", nil, &out, &out)
	if code != 0 {
		t.Fatalf("code=%d out=%q", code, out.String())
	}
	if !strings.Contains(out.String(), "hello-world") {
		t.Fatalf("out=%q", out.String())
	}
}

func TestRunAndOr(t *testing.T) {
	var out bytes.Buffer
	code := RunLine("false && echo should-not", nil, &out, &out)
	if code != 1 {
		t.Fatalf("false && : code=%d", code)
	}
	if strings.Contains(out.String(), "should-not") {
		t.Fatalf("&& ran right side: %q", out.String())
	}

	out.Reset()
	code = RunLine("false || echo fallback", nil, &out, &out)
	if code != 0 {
		t.Fatalf("false || : code=%d out=%q", code, out.String())
	}
	if !strings.Contains(out.String(), "fallback") {
		t.Fatalf("|| missing fallback: %q", out.String())
	}
}

func TestWindowsPathTokenize(t *testing.T) {
	toks, err := Tokenize(`cat C:\Users\User\file.txt | grep x`)
	if err != nil {
		t.Fatal(err)
	}
	if toks[1] != `C:\Users\User\file.txt` {
		t.Fatalf("path broken: %#v", toks)
	}
}
