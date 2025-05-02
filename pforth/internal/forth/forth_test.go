package forth

import (
	"bytes"
	"strings"
	"testing"
)

func newTestForth(input string) (*Forth, *bytes.Buffer) {
	var buf bytes.Buffer
	f := New(strings.NewReader(input), &buf)
	return f, &buf
}

func TestHexPrefix(t *testing.T) {
	f, _ := newTestForth("")
	f.InterpretLine("$FF")
	got := f.DSPop()
	if got != 255 {
		t.Errorf("$FF: expected 255, got %d", got)
	}
}

func TestBinaryPrefix(t *testing.T) {
	f, _ := newTestForth("")
	f.InterpretLine("%1010")
	got := f.DSPop()
	if got != 10 {
		t.Errorf("%%1010: expected 10, got %d", got)
	}
}

func TestSTATEResetOnError(t *testing.T) {
	f, out := newTestForth("")
	f.InterpretLine("UNDEFINED")
	if !strings.Contains(out.String(), "?") {
		t.Errorf("Expected error on undefined word")
	}
	if f.State {
		t.Errorf("STATE should be false after undefined word error")
	}
}
