package forth

import (
	"bufio"
	"io"
	"strings"
)

type Forth struct {
	DS []Cell
	RS []Cell

	Words       []*Word
	wordsByName map[string]int

	Body []Cell
	IP   int

	State bool
	Base  int

	compileName string
	compileList []Cell
	compileWord *Word

	Memory []byte
	DP     int

	TIB []byte
	IN  int

	Input  io.Reader
	Output io.Writer
	bufIn  *bufio.Reader

	Running bool
	Latest  int

	Loaded      map[string]bool
	Interrupted int32
}

func New(input io.Reader, output io.Writer) *Forth {
	f := &Forth{
		DS:          make([]Cell, 0, 64),
		RS:          make([]Cell, 0, 64),
		Words:       make([]*Word, 0, 128),
		wordsByName: make(map[string]int),
		Input:       input,
		Output:      output,
		bufIn:       bufio.NewReader(input),
		Running:     true,
		Memory:      make([]byte, 65536),
		DP:          1024,
		Loaded:      make(map[string]bool),
	}
	initPrimitives(f)
	return f
}

func (f *Forth) DSPush(v Cell) {
	f.DS = append(f.DS, v)
}

func (f *Forth) DSPop() Cell {
	if len(f.DS) == 0 {
		forthError("DATA STACK UNDERFLOW")
	}
	v := f.DS[len(f.DS)-1]
	f.DS = f.DS[:len(f.DS)-1]
	return v
}

func (f *Forth) DSPeek() Cell {
	if len(f.DS) == 0 {
		forthError("DATA STACK UNDERFLOW")
	}
	return f.DS[len(f.DS)-1]
}

func (f *Forth) RSPush(v Cell) {
	f.RS = append(f.RS, v)
}

func (f *Forth) RSPop() Cell {
	if len(f.RS) == 0 {
		forthError("RETURN STACK UNDERFLOW")
	}
	v := f.RS[len(f.RS)-1]
	f.RS = f.RS[:len(f.RS)-1]
	return v
}

func (f *Forth) RSPeek() Cell {
	if len(f.RS) == 0 {
		forthError("RETURN STACK UNDERFLOW")
	}
	return f.RS[len(f.RS)-1]
}

func (f *Forth) Lookup(name string) (*Word, bool) {
	idx, ok := f.wordsByName[name]
	if !ok {
		return nil, false
	}
	return f.Words[idx], true
}

func (f *Forth) FindXT(name string) (int, bool) {
	idx, ok := f.wordsByName[name]
	return idx, ok
}

func (f *Forth) DefineWord(name string, wtype WordType, code XTCode, body []Cell, immediate bool) *Word {
	name = strings.ToUpper(name)
	idx := len(f.Words)
	w := &Word{
		Name:      name,
		Immediate: immediate,
		Type:      wtype,
		Code:      code,
		Body:      body,
		PFA:       f.DP,
	}
	f.Words = append(f.Words, w)
	f.wordsByName[name] = idx
	f.Latest = idx
	return w
}

func (f *Forth) DefinePrimitive(name string, code XTCode) *Word {
	return f.DefineWord(name, WordPrimitive, code, nil, false)
}

func (f *Forth) DefineImmediate(name string, code XTCode) *Word {
	return f.DefineWord(name, WordPrimitive, code, nil, true)
}

func (f *Forth) LatestWord() *Word {
	if f.Latest < 0 || f.Latest >= len(f.Words) {
		return nil
	}
	return f.Words[f.Latest]
}

func (f *Forth) Fetch(addr int) Cell {
	if addr < 0 || addr+1 >= len(f.Memory) {
		forthError("FETCH ADDR OUT OF RANGE: %d", addr)
	}
	return Cell(int(f.Memory[addr]) | int(f.Memory[addr+1])<<8)
}

func (f *Forth) Store(addr int, val Cell) {
	if addr < 0 || addr+1 >= len(f.Memory) {
		forthError("STORE ADDR OUT OF RANGE: %d", addr)
	}
	f.Memory[addr] = byte(val)
	f.Memory[addr+1] = byte(val >> 8)
}

func (f *Forth) CFetch(addr int) byte {
	if addr < 0 || addr >= len(f.Memory) {
		forthError("C@ ADDR OUT OF RANGE: %d", addr)
	}
	return f.Memory[addr]
}

func (f *Forth) CStore(addr int, val byte) {
	if addr < 0 || addr >= len(f.Memory) {
		forthError("C! ADDR OUT OF RANGE: %d", addr)
	}
	f.Memory[addr] = val
}

func (f *Forth) Comma(val Cell) {
	f.Store(f.DP, val)
	f.DP += 2
}

func (f *Forth) CComma(val byte) {
	f.CStore(f.DP, val)
	f.DP++
}

func (f *Forth) Allot(n int) {
	f.DP += n
}

func (f *Forth) Here() int {
	return f.DP
}
