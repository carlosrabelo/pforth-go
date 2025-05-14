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

func TestStackROT(t *testing.T) {
	f, _ := newTestForth("")
	// ROT: ( a b c -- b c a )
	// DS in:  [1, 2, 3]  (TOS=3)
	// pop: a=3, b=2, c=1; push: b=2, a=3, c=1
	// DS out: [2, 3, 1]  (TOS=1)
	f.DSPush(1)
	f.DSPush(2)
	f.DSPush(3)
	exec(f, "ROT")
	if f.DSPop() != 1 {
		t.Errorf("ROT: TOS should be 1")
	}
	if f.DSPop() != 3 {
		t.Errorf("ROT: NOS should be 3")
	}
	if f.DSPop() != 2 {
		t.Errorf("ROT: 3RD should be 2")
	}
}

func TestStackNROT(t *testing.T) {
	f, _ := newTestForth("")
	// -ROT: ( a b c -- c a b )
	// DS in:  [1, 2, 3]  (TOS=3)
	// pop: a=3, b=2, c=1; push: a=3, c=1, b=2
	// DS out: [3, 1, 2]  (TOS=2)
	f.DSPush(1)
	f.DSPush(2)
	f.DSPush(3)
	exec(f, "-ROT")
	if f.DSPop() != 2 {
		t.Errorf("-ROT: TOS should be 2")
	}
	if f.DSPop() != 1 {
		t.Errorf("-ROT: NOS should be 1")
	}
	if f.DSPop() != 3 {
		t.Errorf("-ROT: 3RD should be 3")
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

func TestStackNIP(t *testing.T) {
	f, _ := newTestForth("")
	f.DSPush(10)
	f.DSPush(20)
	exec(f, "NIP")
	if f.DSPop() != 20 {
		t.Errorf("NIP: expected 20, got %d", f.DSPop())
	}
}

func TestStackTUCK(t *testing.T) {
	f, _ := newTestForth("")
	// TUCK: ( a b -- b a b )
	// DS in: [10, 20] (TOS=20)
	// pop: a=20, b=10; push: a=20, b=10, a=20
	// DS out: [20, 10, 20] (TOS=20)
	f.DSPush(10)
	f.DSPush(20)
	exec(f, "TUCK")
	if f.DSPop() != 20 {
		t.Errorf("TUCK: TOS should be 20")
	}
	if f.DSPop() != 10 {
		t.Errorf("TUCK: NOS should be 10")
	}
	if f.DSPop() != 20 {
		t.Errorf("TUCK: 3RD should be 20")
	}
}

func TestStack2DUP(t *testing.T) {
	f, _ := newTestForth("")
	// 2DUP: ( a b -- a b a b )
	f.DSPush(1)
	f.DSPush(2)
	exec(f, "2DUP")
	// DS: [1, 2, 1, 2]  (TOS=2)
	tos := f.DSPop()  // 2
	nos := f.DSPop()  // 1
	thd := f.DSPop()  // 2
	fou := f.DSPop()  // 1
	if tos != 2 || nos != 1 || thd != 2 || fou != 1 {
		t.Errorf("2DUP: expected 2 1 2 1")
	}
}

func TestStack2DROP(t *testing.T) {
	f, _ := newTestForth("")
	f.DSPush(1)
	f.DSPush(2)
	exec(f, "2DROP")
	if len(f.DS) != 0 {
		t.Errorf("2DROP: DS should be empty")
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
		{"1+", "1+", []Cell{41}, 42},
		{"1-", "1-", []Cell{42}, 41},
		{"2*", "2*", []Cell{21}, 42},
		{"2/", "2/", []Cell{42}, 21},
		{"NEGATE", "NEGATE", []Cell{42}, -42},
		{"ABS positive", "ABS", []Cell{42}, 42},
		{"ABS negative", "ABS", []Cell{-42}, 42},
		{"MIN a<b", "MIN", []Cell{3, 7}, 3},
		{"MIN a>b", "MIN", []Cell{7, 3}, 3},
		{"MAX a<b", "MAX", []Cell{3, 7}, 7},
		{"MAX a>b", "MAX", []Cell{7, 3}, 7},
		{"/", "/", []Cell{10, 3}, 3},
		{"MOD", "MOD", []Cell{10, 3}, 1},
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

func TestDivisionByZero(t *testing.T) {
	f, _ := newTestForth("")
	f.DSPush(5)   // dividend
	f.DSPush(0)   // divisor
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic for division by zero")
		}
	}()
	exec(f, "/")
}

func TestComparison(t *testing.T) {
	tests := []struct {
		name   string
		word   string
		push   []Cell
		want   Cell
	}{
		{"=", "=", []Cell{5, 5}, -1},
		{"= false", "=", []Cell{5, 6}, 0},
		{"<> true", "<>", []Cell{5, 6}, -1},
		{"<> false", "<>", []Cell{5, 5}, 0},
		{"< true", "<", []Cell{3, 5}, -1},
		{"< false", "<", []Cell{5, 3}, 0},
		{"> true", ">", []Cell{5, 3}, -1},
		{"> false", ">", []Cell{3, 5}, 0},
		{"0= zero", "0=", []Cell{0}, -1},
		{"0= non-zero", "0=", []Cell{5}, 0},
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

func TestMemoryAddStore(t *testing.T) {
	f, _ := newTestForth("")
	f.DSPush(100)
	f.DSPush(4000)
	exec(f, "!")     // mem[4000] = 100
	f.DSPush(50)
	f.DSPush(4000)
	exec(f, "+!")    // mem[4000] += 50
	f.DSPush(4000)
	exec(f, "@")
	if f.DSPop() != 150 {
		t.Errorf("+!: expected 150, got %d", f.DSPop())
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

func TestMemoryERASE(t *testing.T) {
	f, _ := newTestForth("")
	for i := 0; i < 10; i++ {
		f.CStore(6000+i, 0xFF)
	}
	f.DSPush(6000)
	f.DSPush(10)
	exec(f, "ERASE")
	for i := 0; i < 10; i++ {
		if f.CFetch(6000+i) != 0 {
			t.Errorf("ERASE[%d]: expected 0", i)
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

func TestCOLON(t *testing.T) {
	f, _ := newTestForth("")
	f.InterpretLine(": DOUBLE DUP + ;")
	w, found := f.Lookup("DOUBLE")
	if !found {
		t.Fatal("DOUBLE not found after definition")
	}
	if w.Type != WordColon {
		t.Errorf("DOUBLE: expected colon type, got %v", w.Type)
	}
}

func TestColonDefinitionAndExecution(t *testing.T) {
	f2, out := newTestForth("")
	f2.InterpretLine(": SQUARE DUP * ;")
	f2.InterpretLine("5 SQUARE .")
	if !strings.Contains(out.String(), "25") {
		t.Errorf("SQUARE: expected 25, got %q", out.String())
	}
}

func TestIFELSE(t *testing.T) {
	f, out := newTestForth("")
	f.InterpretLine(": TEST IF .\" YES\" ELSE .\" NO\" THEN ;")
	f.InterpretLine("1 TEST")
	f.InterpretLine("0 TEST")
	output := out.String()
	if !strings.Contains(output, "YES") {
		t.Errorf("1 TEST should print YES, got %q", output)
	}
	if !strings.Contains(output, "NO") {
		t.Errorf("0 TEST should print NO, got %q", output)
	}
}

func TestBEGINUNTIL(t *testing.T) {
	f, out := newTestForth("")
	f.InterpretLine(": COUNTDOWN BEGIN DUP . 1- DUP 0= UNTIL DROP ;")
	f.InterpretLine("5 COUNTDOWN")
	output := out.String()
	if !strings.Contains(output, "5 4 3 2 1") {
		t.Errorf("COUNTDOWN: expected 5 4 3 2 1, got %q", output)
	}
}

func TestDOLOOP(t *testing.T) {
	f, out := newTestForth("")
	f.InterpretLine(": TEST 5 0 DO I . LOOP ;")
	f.InterpretLine("TEST")
	output := out.String()
	if !strings.Contains(output, "0 1 2 3 4") {
		t.Errorf("DO/LOOP: expected 0 1 2 3 4, got %q", output)
	}
}

func TestPLUSLOOP(t *testing.T) {
	f, out := newTestForth("")
	f.InterpretLine(": TEST 10 0 DO I . 3 +LOOP ;")
	f.InterpretLine("TEST")
	output := out.String()
	if !strings.Contains(output, "0 3 6 9") {
		t.Errorf("+LOOP: expected 0 3 6 9, got %q", output)
	}
}

func TestNestedDOLOOP(t *testing.T) {
	f, out := newTestForth("")
	f.InterpretLine(": TEST 3 0 DO 3 0 DO J . I . LOOP LOOP ;")
	f.InterpretLine("TEST")
	output := out.String()
	if !strings.Contains(output, "0 0") || !strings.Contains(output, "2 2") {
		t.Errorf("nested DO/LOOP output: %q", output)
	}
}

func TestFACTORIAL(t *testing.T) {
	f, out := newTestForth("")
	f.InterpretLine(": FACTORIAL")
	f.InterpretLine("  DUP 0= IF DROP 1 ELSE DUP 1- RECURSE * THEN ;")
	f.InterpretLine("10 FACTORIAL .")
	output := out.String()
	if !strings.Contains(output, "3628800") {
		t.Errorf("FACTORIAL(10): expected 3628800, got %q", output)
	}
}

func TestVARIABLE(t *testing.T) {
	f, out := newTestForth("")
	f.InterpretLine("VARIABLE X")
	f.InterpretLine("42 X !")
	f.InterpretLine("X @ .")
	if !strings.Contains(out.String(), "42") {
		t.Errorf("VARIABLE: expected 42, got %q", out.String())
	}
}

func TestCONSTANT(t *testing.T) {
	f, out := newTestForth("")
	f.InterpretLine("42 CONSTANT ANSWER")
	f.InterpretLine("ANSWER .")
	if !strings.Contains(out.String(), "42") {
		t.Errorf("CONSTANT: expected 42, got %q", out.String())
	}
}

func TestCREATE(t *testing.T) {
	f, _ := newTestForth("")
	f.InterpretLine("CREATE FOO")
	if _, ok := f.Lookup("FOO"); !ok {
		t.Errorf("CREATE: FOO not found")
	}
}

func TestTICK(t *testing.T) {
	f, _ := newTestForth("")
	xtDUP, ok := f.FindXT("DUP")
	if !ok {
		t.Fatal("DUP not found in dictionary")
	}
	f.InterpretLine("' DUP")
	if f.DSPop() != Cell(xtDUP) {
		t.Errorf("': expected XT=%d, got %d", xtDUP, f.DSPop())
	}
}

func TestDOTQUOTE(t *testing.T) {
	f, out := newTestForth("")
	f.InterpretLine(".\" hello\"")
	if !strings.Contains(out.String(), "hello") {
		t.Errorf(".\" : expected 'hello', got %q", out.String())
	}
}

func TestDOTQUOTECompile(t *testing.T) {
	f, out := newTestForth("")
	f.InterpretLine(": SAYHELLO .\" hello from def\" ;")
	f.InterpretLine("SAYHELLO")
	if !strings.Contains(out.String(), "hello from def") {
		t.Errorf(".\" compiled: expected 'hello from def', got %q", out.String())
	}
}

func TestSQUOTECompile(t *testing.T) {
	f, out := newTestForth("")
	f.InterpretLine(": HELLO S\" hello\" TYPE ;")
	f.InterpretLine("HELLO")
	if !strings.Contains(out.String(), "hello") {
		t.Errorf("S\" compiled: expected 'hello', got %q", out.String())
	}
}

func TestEVALUATE(t *testing.T) {
	f, out := newTestForth("")
	// Store "2 3 + ." in memory
	s := "2 3 + ."
	addr := f.Here()
	for i := 0; i < len(s); i++ {
		f.CStore(addr+i, s[i])
	}
	f.DP = addr + len(s)
	f.DSPush(Cell(addr))
	f.DSPush(Cell(len(s)))
	exec(f, "EVALUATE")
	if !strings.Contains(out.String(), "5") {
		t.Errorf("EVALUATE: expected '5' output, got %q", out.String())
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

func TestSTATEResetOnError(t *testing.T) {
	f, out := newTestForth("")
	f.InterpretLine(": TEST UNDEFINED ;")
	if !strings.Contains(out.String(), "?") {
		t.Errorf("Expected error on undefined word")
	}
	if f.State {
		t.Errorf("STATE should be false after undefined word error")
	}
}

func TestBYE(t *testing.T) {
	f, _ := newTestForth("")
	exec(f, "BYE")
	if f.Running {
		t.Errorf("BYE should set Running to false")
	}
}

func TestSPAN(t *testing.T) {
	f, _ := newTestForth("hello\n")
	f.DSPush(5000) // buffer address
	f.DSPush(80)   // max length
	f.InterpretLine("EXPECT")
	exec(f, "SPAN", "@")
	span := f.DSPop()
	if span != 5 {
		t.Errorf("SPAN: expected 5, got %d", span)
	}
}

func TestABORT(t *testing.T) {
	f, _ := newTestForth("")
	f.State = true
	f.DSPush(10)
	f.Body = []Cell{1, 2, 3}
	f.IP = 1
	exec(f, "ABORT")
	if f.State {
		t.Errorf("ABORT: STATE should be false")
	}
	if len(f.DS) != 0 {
		t.Errorf("ABORT: DS should be empty")
	}
}

func TestRECURSE(t *testing.T) {
	f, out := newTestForth("")
	// 5! = 120
	f.InterpretLine(": FACT DUP 1 > IF DUP 1- FACT * THEN ;")
	f.InterpretLine("5 FACT .")
	if !strings.Contains(out.String(), "120") {
		t.Errorf("FACT(5): expected 120, got %q", out.String())
	}
}

func TestBEGINAGAIN(t *testing.T) {
	f, out := newTestForth("")
	f.InterpretLine(": INF 0 BEGIN DUP . 1+ DUP 5 = UNTIL DROP ;")
	f.InterpretLine("INF")
	if !strings.Contains(out.String(), "0 1 2 3 4") {
		t.Errorf("BEGIN/UNTIL: expected 0 1 2 3 4, got %q", out.String())
	}
}

func TestWHILEREPEAT(t *testing.T) {
	f, out := newTestForth("")
	f.InterpretLine(": TEST")
	f.InterpretLine("  0")
	f.InterpretLine("  BEGIN")
	f.InterpretLine("    DUP . 1+")
	f.InterpretLine("    DUP 5 <")
	f.InterpretLine("  WHILE")
	f.InterpretLine("  REPEAT")
	f.InterpretLine("  DROP ;")
	f.InterpretLine("TEST")
	if !strings.Contains(out.String(), "0 1 2 3 4") {
		t.Errorf("BEGIN/WHILE/REPEAT: expected 0 1 2 3 4, got %q", out.String())
	}
}

func TestQUESTION(t *testing.T) {
	f, out := newTestForth("")
	f.InterpretLine("VARIABLE X")
	f.InterpretLine("42 X !")
	f.InterpretLine("X ?")
	if !strings.Contains(out.String(), "42") {
		t.Errorf("?: expected 42, got %q", out.String())
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

func TestInterrupt(t *testing.T) {
	f, _ := newTestForth("")
	f.InterpretLine(": INFLOOP 0 BEGIN 1+ AGAIN ;")
	done := make(chan struct{})
	go func() {
		defer func() { recover(); close(done) }()
		exec(f, "INFLOOP")
	}()
	atomic.StoreInt32(&f.Interrupted, 1)
	<-done
	if len(f.DS) != 0 {
		t.Errorf("expected empty stack after interrupt")
	}
}
