package forth

import (
	"bytes"
	"strings"
	"sync/atomic"
	"testing"
)

func newTestForth(input string) (*Forth, *bytes.Buffer) {
	var buf bytes.Buffer
	f := New(strings.NewReader(input), &buf)
	return f, &buf
}

func exec(f *Forth, words ...string) {
	for _, w := range words {
		if xt, ok := f.FindXT(w); ok {
			f.ExecuteWord(xt)
		}
	}
}

func TestStackDUP(t *testing.T) {
	f, _ := newTestForth("")
	f.DSPush(42)
	f.DSPush(0) // placeholder, we'll check DUP
	f.DSPop()
	exec(f, "DUP")
	if f.DSPop() != 42 || f.DSPop() != 42 {
		t.Errorf("DUP: expected 42 42")
	}
}

func TestStackDROP(t *testing.T) {
	f, _ := newTestForth("")
	f.DSPush(10)
	f.DSPush(20)
	exec(f, "DROP")
	if f.DSPop() != 10 {
		t.Errorf("DROP: expected 10, got %d", f.DSPop())
	}
}

func TestStackSWAP(t *testing.T) {
	f, _ := newTestForth("")
	f.DSPush(10)
	f.DSPush(20)
	exec(f, "SWAP")
	// SWAP: ( 10 20 -- 20 10 ), so TOS=10
	if f.DSPop() != 10 {
		t.Errorf("SWAP: TOS should be 10")
	}
	if f.DSPop() != 20 {
		t.Errorf("SWAP: NOS should be 20")
	}
}

func TestStackOVER(t *testing.T) {
	f, _ := newTestForth("")
	f.DSPush(10)
	f.DSPush(20)
	exec(f, "OVER")
	// OVER: ( 10 20 -- 10 20 10 ), TOS=10
	if f.DSPop() != 10 {
		t.Errorf("OVER: TOS should be 10")
	}
	if f.DSPop() != 20 {
		t.Errorf("OVER: NOS should be 20")
	}
	if f.DSPop() != 10 {
		t.Errorf("OVER: 3RD should be 10")
	}
}

func TestStackQDUP(t *testing.T) {
	f, _ := newTestForth("")
	f.DSPush(0)
	exec(f, "?DUP")
	if len(f.DS) != 1 {
		t.Errorf("?DUP(0): expected depth 1, got %d", len(f.DS))
	}
	f.DSPush(5)
	exec(f, "?DUP")
	if len(f.DS) != 3 {
		t.Errorf("?DUP(5): expected depth 3, got %d", len(f.DS))
	}
	if f.DSPop() != 5 {
		t.Errorf("?DUP(5): top should be 5")
	}
}

func TestStackToRFromR(t *testing.T) {
	f, _ := newTestForth("")
	f.DSPush(42)
	exec(f, ">R")
	if len(f.DS) != 0 {
		t.Errorf(">R: DS should be empty, got %d", len(f.DS))
	}
	if len(f.RS) != 1 {
		t.Errorf(">R: RS should have 1 item, got %d", len(f.RS))
	}
	exec(f, "R>")
	if f.DSPop() != 42 {
		t.Errorf("R>: expected 42")
	}
}

func TestStackRFetch(t *testing.T) {
	f, _ := newTestForth("")
	f.DSPush(99)
	exec(f, ">R", "R@")
	if f.DSPop() != 99 {
		t.Errorf("R@: expected 99")
	}
	if len(f.RS) != 1 {
		t.Errorf("R@: RS should still have 1 item")
	}
}

func TestStack2SWAP(t *testing.T) {
	f, _ := newTestForth("")
	// 2SWAP: ( a b c d -- c d a b )
	f.DSPush(1)
	f.DSPush(2)
	f.DSPush(3)
	f.DSPush(4)
	exec(f, "2SWAP")
	// DS: [1, 2, 3, 4] -> after: [3, 4, 1, 2]  (TOS=2)
	tos := f.DSPop()  // 2
	nos := f.DSPop()  // 1
	thd := f.DSPop()  // 4
	fou := f.DSPop()  // 3
	if tos != 2 || nos != 1 || thd != 4 || fou != 3 {
		t.Errorf("2SWAP: expected 2 1 4 3")
	}
}

func TestStack2OVER(t *testing.T) {
	f, _ := newTestForth("")
	f.DSPush(1)
	f.DSPush(2)
	f.DSPush(3)
	f.DSPush(4)
	exec(f, "2OVER")
	// Current impl: ( a b c d -- a b c d c d ) with TOS=d
	// DS out: [3, 4, 1, 2, 3, 4] (TOS=4)
	popOrder := make([]Cell, 6)
	for i := 0; i < 6; i++ {
		popOrder[i] = f.DSPop()
	}
	expected := []Cell{4, 3, 2, 1, 4, 3}
	for i, v := range popOrder {
		if v != expected[i] {
			t.Errorf("2OVER pop[%d]: expected %d, got %d", i, expected[i], v)
		}
	}
}

func TestDEPTH(t *testing.T) {
	f, _ := newTestForth("")
	if len(f.DS) != 0 {
		// drain stack
		for len(f.DS) > 0 {
			f.DSPop()
		}
	}
	f.DSPush(10)
	f.DSPush(20)
	exec(f, "DEPTH")
	if f.DSPop() != 2 {
		t.Errorf("DEPTH: expected 2, got %d", f.DSPop())
	}
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

func TestInterrupt(t *testing.T) {
	f, _ := newTestForth("")
	atomic.StoreInt32(&f.Interrupted, 1)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected interrupt")
		}
	}()
	f.InterpretLine("0")
}
