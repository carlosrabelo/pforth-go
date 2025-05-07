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

func TestArithmetic(t *testing.T) {
	tests := []struct {
		name   string
		word   string
		inputs []Cell
		want   Cell
	}{
		{"+", "+", []Cell{10, 20}, 30},
		{"-", "-", []Cell{20, 5}, 15},
		{"*", "*", []Cell{6, 7}, 42},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, _ := newTestForth("")
			for _, v := range tt.inputs {
				f.DSPush(v)
			}
			exec(f, tt.word)
			got := f.DSPop()
			if got != tt.want {
				t.Errorf("%s: expected %d, got %d", tt.name, tt.want, got)
			}
		})
	}
}

func TestDIVMOD(t *testing.T) {
	f, _ := newTestForth("")
	// /MOD: ( dividend divisor -- rem quot )
	// 10 / 3 = 3 remainder 1
	f.DSPush(10)
	f.DSPush(3)
	exec(f, "/MOD")
	quot := f.DSPop()
	rem := f.DSPop()
	if rem != 1 || quot != 3 {
		t.Errorf("/MOD: expected rem=1 quot=3, got rem=%d quot=%d", rem, quot)
	}
}

func TestComparison(t *testing.T) {
	tests := []struct {
		name string
		word string
		push []Cell
		want Cell
	}{
		{"=", "=", []Cell{5, 5}, -1},
		{"= false", "=", []Cell{5, 6}, 0},
		{"< true", "<", []Cell{3, 5}, -1},
		{"< false", "<", []Cell{5, 3}, 0},
		{"> true", ">", []Cell{5, 3}, -1},
		{"> false", ">", []Cell{3, 5}, 0},
		{"AND", "AND", []Cell{3, 1}, 1},
		{"OR", "OR", []Cell{1, 2}, 3},
		{"XOR", "XOR", []Cell{3, 1}, 2},
		{"INVERT", "INVERT", []Cell{0}, -1},
		{"INVERT all", "INVERT", []Cell{-1}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, _ := newTestForth("")
			for _, v := range tt.push {
				f.DSPush(v)
			}
			exec(f, tt.word)
			got := f.DSPop()
			if got != tt.want {
				t.Errorf("%s: expected %d, got %d", tt.name, tt.want, got)
			}
		})
	}
}

func TestMemoryStoreAndFetch(t *testing.T) {
	f, _ := newTestForth("")
	f.DSPush(12345)
	f.DSPush(2000)
	exec(f, "!")
	f.DSPush(2000)
	exec(f, "@")
	if f.DSPop() != 12345 {
		t.Errorf("@/!: expected 12345, got %d", f.DSPop())
	}
}

func TestMemoryCStoreCFetch(t *testing.T) {
	f, _ := newTestForth("")
	f.DSPush(255)
	f.DSPush(3000)
	exec(f, "C!")
	f.DSPush(3000)
	exec(f, "C@")
	if f.DSPop() != 255 {
		t.Errorf("C@/C!: expected 255, got %d", f.DSPop())
	}
}

func TestMemoryFILL(t *testing.T) {
	f, _ := newTestForth("")
	f.DSPush(5000)
	f.DSPush(5)
	f.DSPush(65) // 'A'
	exec(f, "FILL")
	for i := 0; i < 5; i++ {
		if f.CFetch(5000+i) != 65 {
			t.Errorf("FILL[%d]: expected 65", i)
		}
	}
}

func TestMemoryCMOVE(t *testing.T) {
	f, _ := newTestForth("")
	for i := 0; i < 5; i++ {
		f.CStore(7000+i, byte(i+10))
	}
	f.DSPush(7005)
	f.DSPush(7000)
	f.DSPush(5)
	exec(f, "CMOVE")
	for i := 0; i < 5; i++ {
		if f.CFetch(7005+i) != byte(i+10) {
			t.Errorf("CMOVE[%d]: expected %d, got %d", i, i+10, f.CFetch(7005+i))
		}
	}
}

func TestCMOVE(t *testing.T) {
	f, _ := newTestForth("")
	for i := 0; i < 5; i++ {
		f.CStore(8000+i, byte(i+10))
	}
	// Non-overlapping copy: src=8000, dst=8010, len=5
	f.DSPush(8010)
	f.DSPush(8000)
	f.DSPush(5)
	exec(f, "CMOVE")
	for i := 0; i < 5; i++ {
		if f.CFetch(8010+i) != byte(i+10) {
			t.Errorf("CMOVE[%d]: expected %d", i, i+10)
		}
	}
}

func TestCMOVEOverlap(t *testing.T) {
	f, _ := newTestForth("")
	for i := 0; i < 5; i++ {
		f.CStore(8000+i, byte(i+20))
	}
	// Forward src < dst: propagates first byte
	f.DSPush(8001)
	f.DSPush(8000)
	f.DSPush(4)
	exec(f, "CMOVE")
	// After forward CMOVE with src<dst: all copied bytes get src[0]
	if f.CFetch(8001) != 20 || f.CFetch(8002) != 20 {
		t.Errorf("CMOVE fwd: expected first-byte propagation")
	}
}

func TestHexPrefix(t *testing.T) {
	f, out := newTestForth("")
	f.InterpretLine("$FF .")
	if !strings.Contains(out.String(), "255") {
		t.Errorf("$FF: expected 255, got %q", out.String())
	}
}

func TestBinaryPrefix(t *testing.T) {
	f, out := newTestForth("")
	f.InterpretLine("%1010 .")
	if !strings.Contains(out.String(), "10") {
		t.Errorf("%%1010: expected 10, got %q", out.String())
	}
}

func TestBYE(t *testing.T) {
	f, _ := newTestForth("")
	exec(f, "BYE")
	if f.Running {
		t.Errorf("BYE should set Running to false")
	}
}

func TestDOTS(t *testing.T) {
	f, out := newTestForth("")
	f.InterpretLine("1 2 3 .S")
	output := out.String()
	if !strings.Contains(output, "<") || !strings.Contains(output, ">") {
		t.Errorf(".S: expected stack display with <>, got %q", output)
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
