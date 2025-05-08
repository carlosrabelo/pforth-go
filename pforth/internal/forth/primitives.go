package forth

import (
	"strconv"
	"strings"
	"unicode"
)

func initPrimitives(f *Forth) {
	f.DefinePrimitive("DUP", dup)
	f.DefinePrimitive("DROP", drop)
	f.DefinePrimitive("SWAP", swap)
	f.DefinePrimitive("OVER", over)
	f.DefinePrimitive("?DUP", qdup)
	f.DefinePrimitive(">R", tor)
	f.DefinePrimitive("R>", fromr)
	f.DefinePrimitive("R@", rfetch)
	f.DefinePrimitive("2SWAP", twoswap)
	f.DefinePrimitive("2OVER", twover)
	f.DefinePrimitive("DEPTH", depth)

	f.DefinePrimitive("@", fetch)
	f.DefinePrimitive("!", store)
	f.DefinePrimitive("C@", cfetch)
	f.DefinePrimitive("C!", cstore)
	f.DefinePrimitive("FILL", fill)
	f.DefinePrimitive("CMOVE", cmove)
	f.DefinePrimitive("CMOVE>", cmover)
	f.DefinePrimitive("HERE", here)
	f.DefinePrimitive("PAD", pad)
	f.DefinePrimitive("ALLOT", allot)
	f.DefinePrimitive(",", comma)
	f.DefinePrimitive("C,", ccomma)

	f.DefinePrimitive("+", plus)
	f.DefinePrimitive("-", minus)
	f.DefinePrimitive("*", star)
	f.DefinePrimitive("/MOD", slashmod)
	f.DefinePrimitive("U<=", ulteq)
	f.DefinePrimitive("U>=", ugteq)
	f.DefinePrimitive("AND", andop)
	f.DefinePrimitive("OR", orop)
	f.DefinePrimitive("XOR", xorop)
	f.DefinePrimitive("INVERT", invert)
	f.DefinePrimitive("=", eq)
	f.DefinePrimitive("<", lt)
	f.DefinePrimitive(">", gt)
	f.DefinePrimitive("U<", ult)
	f.DefinePrimitive("U>", ugt)

	f.DefinePrimitive("KEY", key)
	f.DefinePrimitive("EMIT", emit)
	f.DefinePrimitive("CR", cr)
	f.DefinePrimitive("TYPE", typeof)
	f.DefinePrimitive(".", dot)
	f.DefinePrimitive("U.", udot)
	f.DefinePrimitive("BYE", bye)
	f.DefinePrimitive(".S", dots)

	f.DefinePrimitive("WORD", word)
	f.DefinePrimitive("FIND", find)
	f.DefinePrimitive("EXECUTE", execute)
	f.DefinePrimitive("INTERPRET", interpret)
	f.DefinePrimitive("QUIT", quit)
	f.DefinePrimitive("ABORT", abortw)
	f.DefinePrimitive("EXPECT", expect)
	f.DefinePrimitive("EVALUATE", evaluate)
	f.DefinePrimitive("BASE", base)
	f.DefinePrimitive("DP", dpaddr)
	f.DefinePrimitive("SPAN", span)
	f.DefinePrimitive("SOURCE", source)

	f.DefinePrimitive("LIT", lit)
	f.DefinePrimitive("EXIT", exitword)

	f.DefinePrimitive(":", colon)
	f.DefineImmediate(";", semicolon)
	f.DefineImmediate("[", lbracket)
	f.DefineImmediate("]", rbracket)
	f.DefinePrimitive("STATE", state)
	f.DefinePrimitive("CREATE", create)
	f.DefinePrimitive("CONSTANT", constant)
	f.DefinePrimitive("VARIABLE", variable)
	f.DefineImmediate("IMMEDIATE", immediate)
	f.DefineImmediate("[COMPILE]", bracketcompile)
	f.DefineImmediate("LITERAL", literal)
	f.DefineImmediate(".\"", dotstring)
	f.DefineImmediate("S\"", sstring)
	f.DefineImmediate("(", parenword)
	f.DefinePrimitive(">BODY", tobody)
	f.DefinePrimitive("'", tick)
	f.DefineImmediate("RECURSE", recurse)

	f.DefinePrimitive("BRANCH", branch)
	f.DefinePrimitive("0BRANCH", zerobranch)

	f.DefineImmediate("IF", ifword)
	f.DefineImmediate("THEN", thenword)
	f.DefineImmediate("ELSE", elseword)
	f.DefineImmediate("BEGIN", beginword)
	f.DefineImmediate("UNTIL", untilword)
	f.DefineImmediate("AGAIN", againword)
	f.DefineImmediate("WHILE", whileword)
	f.DefineImmediate("REPEAT", repeatword)
}

var dup = func(f *Forth) {
	v := f.DSPeek()
	f.DSPush(v)
}

var drop = func(f *Forth) {
	f.DSPop()
}

var swap = func(f *Forth) {
	a := f.DSPop()
	b := f.DSPop()
	f.DSPush(a)
	f.DSPush(b)
}

var over = func(f *Forth) {
	a := f.DSPop()
	b := f.DSPeek()
	f.DSPush(a)
	f.DSPush(b)
}

var qdup = func(f *Forth) {
	v := f.DSPeek()
	if v != 0 {
		f.DSPush(v)
	}
}

var tor = func(f *Forth) {
	v := f.DSPop()
	f.RSPush(v)
}

var fromr = func(f *Forth) {
	v := f.RSPop()
	f.DSPush(v)
}

var rfetch = func(f *Forth) {
	v := f.RSPeek()
	f.DSPush(v)
}

var twoswap = func(f *Forth) {
	a := f.DSPop()
	b := f.DSPop()
	c := f.DSPop()
	d := f.DSPop()
	f.DSPush(b)
	f.DSPush(a)
	f.DSPush(d)
	f.DSPush(c)
}

var twover = func(f *Forth) {
	a := f.DSPop()
	b := f.DSPop()
	c := f.DSPop()
	d := f.DSPop()
	f.DSPush(b)
	f.DSPush(a)
	f.DSPush(d)
	f.DSPush(c)
	f.DSPush(b)
	f.DSPush(a)
}

var depth = func(f *Forth) {
	f.DSPush(Cell(len(f.DS)))
}

var fetch = func(f *Forth) {
	addr := int(f.DSPop())
	val := f.Fetch(addr)
	f.DSPush(val)
}

var store = func(f *Forth) {
	addr := int(f.DSPop())
	val := f.DSPop()
	f.Store(addr, val)
}

var cfetch = func(f *Forth) {
	addr := int(f.DSPop())
	val := Cell(f.CFetch(addr))
	f.DSPush(val)
}

var cstore = func(f *Forth) {
	addr := int(f.DSPop())
	val := byte(f.DSPop())
	f.CStore(addr, val)
}

var fill = func(f *Forth) {
	val := byte(f.DSPop())
	length := int(f.DSPop())
	addr := int(f.DSPop())
	for i := 0; i < length; i++ {
		f.CStore(addr+i, val)
	}
}

var cmove = func(f *Forth) {
	length := int(f.DSPop())
	src := int(f.DSPop())
	dst := int(f.DSPop())
	for i := 0; i < length; i++ {
		f.CStore(dst+i, f.CFetch(src+i))
	}
}

var cmover = func(f *Forth) {
	length := int(f.DSPop())
	src := int(f.DSPop())
	dst := int(f.DSPop())
	for i := length - 1; i >= 0; i-- {
		f.CStore(dst+i, f.CFetch(src+i))
	}
}

var here = func(f *Forth) {
	f.DSPush(Cell(f.Here()))
}

var pad = func(f *Forth) {
	f.DSPush(Cell(f.Here() + 68))
}

var allot = func(f *Forth) {
	n := int(f.DSPop())
	f.Allot(n)
}

var comma = func(f *Forth) {
	val := f.DSPop()
	f.Comma(val)
}

var ccomma = func(f *Forth) {
	val := byte(f.DSPop())
	f.CComma(val)
}

var plus = func(f *Forth) {
	a := f.DSPop()
	b := f.DSPop()
	f.DSPush(a + b)
}

var minus = func(f *Forth) {
	a := f.DSPop()
	b := f.DSPop()
	f.DSPush(b - a)
}

var star = func(f *Forth) {
	a := f.DSPop()
	b := f.DSPop()
	f.DSPush(a * b)
}

var slashmod = func(f *Forth) {
	div := f.DSPop()
	num := f.DSPop()
	if div == 0 {
		forthError("DIVISION BY ZERO")
	}
	quot := num / div
	rem := num % div
	f.DSPush(rem)
	f.DSPush(quot)
}

var ulteq = func(f *Forth) {
	a := f.DSPop()
	b := f.DSPop()
	if uint(b) <= uint(a) {
		f.DSPush(-1)
	} else {
		f.DSPush(0)
	}
}

var ugteq = func(f *Forth) {
	a := f.DSPop()
	b := f.DSPop()
	if uint(b) >= uint(a) {
		f.DSPush(-1)
	} else {
		f.DSPush(0)
	}
}

var andop = func(f *Forth) {
	a := f.DSPop()
	b := f.DSPop()
	f.DSPush(a & b)
}

var orop = func(f *Forth) {
	a := f.DSPop()
	b := f.DSPop()
	f.DSPush(a | b)
}

var xorop = func(f *Forth) {
	a := f.DSPop()
	b := f.DSPop()
	f.DSPush(a ^ b)
}

var invert = func(f *Forth) {
	v := f.DSPop()
	f.DSPush(^v)
}

var eq = func(f *Forth) {
	a := f.DSPop()
	b := f.DSPop()
	if b == a {
		f.DSPush(-1)
	} else {
		f.DSPush(0)
	}
}

var lt = func(f *Forth) {
	a := f.DSPop()
	b := f.DSPop()
	if b < a {
		f.DSPush(-1)
	} else {
		f.DSPush(0)
	}
}

var gt = func(f *Forth) {
	a := f.DSPop()
	b := f.DSPop()
	if b > a {
		f.DSPush(-1)
	} else {
		f.DSPush(0)
	}
}

var ult = func(f *Forth) {
	a := f.DSPop()
	b := f.DSPop()
	if uint(b) < uint(a) {
		f.DSPush(-1)
	} else {
		f.DSPush(0)
	}
}

var ugt = func(f *Forth) {
	a := f.DSPop()
	b := f.DSPop()
	if uint(b) > uint(a) {
		f.DSPush(-1)
	} else {
		f.DSPush(0)
	}
}

var key = func(f *Forth) {
	b, err := f.bufIn.ReadByte()
	if err != nil {
		f.DSPush(0)
		return
	}
	f.DSPush(Cell(b))
}

var emit = func(f *Forth) {
	c := byte(f.DSPop())
	f.Emit(c)
}

var cr = func(f *Forth) {
	f.EmitCR()
}

var typeof = func(f *Forth) {
	length := int(f.DSPop())
	addr := int(f.DSPop())
	for i := 0; i < length; i++ {
		f.Emit(f.CFetch(addr + i))
	}
}

var dot = func(f *Forth) {
	n := f.DSPop()
	s := strconv.FormatInt(int64(n), 10)
	f.EmitStr(s)
	f.Emit(' ')
}

var udot = func(f *Forth) {
	n := f.DSPop()
	s := strconv.FormatUint(uint64(n), 10)
	f.EmitStr(s)
	f.Emit(' ')
}

var bye = func(f *Forth) {
	f.Running = false
}

var dots = func(f *Forth) {
	f.EmitStr("<")
	depth := len(f.DS)
	for i := depth - 1; i >= 0; i-- {
		s := strconv.FormatInt(int64(f.DS[i]), 10)
		f.EmitStr(s)
		if i > 0 {
			f.Emit(' ')
		}
	}
	f.EmitStr(">")
}

var word = func(f *Forth) {
	delim := byte(f.DSPop())
	for f.IN < len(f.TIB) && f.TIB[f.IN] == delim {
		f.IN++
	}
	if f.IN >= len(f.TIB) {
		f.CStore(f.Here(), 0)
		f.DSPush(Cell(f.Here()))
		return
	}
	start := f.Here() + 1
	count := 0
	for f.IN < len(f.TIB) &&
		f.TIB[f.IN] != delim &&
		f.TIB[f.IN] != '\r' &&
		f.TIB[f.IN] != '\n' {
		f.CStore(start+count, f.TIB[f.IN])
		count++
		f.IN++
	}
	f.CStore(f.Here(), byte(count))
	f.CStore(start+count, ' ')
	f.DSPush(Cell(f.Here()))
}

var find = func(f *Forth) {
	addr := int(f.DSPop())
	length := int(f.CFetch(addr))
	nameBytes := make([]byte, length)
	for i := 0; i < length; i++ {
		nameBytes[i] = byte(unicode.ToUpper(rune(f.CFetch(addr + 1 + i))))
	}
	name := string(nameBytes)
	if _, ok := f.Lookup(name); !ok {
		f.DSPush(Cell(addr))
		f.DSPush(0)
		return
	}
	xt, _ := f.FindXT(name)
	f.DSPush(Cell(xt))
	f.DSPush(1)
}

var execute = func(f *Forth) {
	xt := int(f.DSPop())
	f.ExecuteWord(xt)
}

var interpret = func(f *Forth) {
	f.interpretLoop()
}

var quit = func(f *Forth) {
	f.QUIT()
}

var abortw = func(f *Forth) {
	f.Abort()
}

var expect = func(f *Forth) {
	maxlen := int(f.DSPop())
	addr := int(f.DSPop())
	line, err := f.ReadLine()
	if err != nil {
		return
	}
	count := len(line)
	if count > maxlen {
		count = maxlen
	}
	for i := 0; i < count; i++ {
		f.CStore(addr+i, line[i])
	}
	f.CStore(addr+count, 0)
	f.Store(1000, Cell(count))
}

var evaluate = func(f *Forth) {
	len := int(f.DSPop())
	addr := int(f.DSPop())
	savedTIB := f.TIB
	savedIN := f.IN
	line := make([]byte, len)
	for i := 0; i < len; i++ {
		line[i] = f.CFetch(addr + i)
	}
	f.TIB = line
	f.IN = 0
	f.interpretLoop()
	f.TIB = savedTIB
	f.IN = savedIN
}

var base = func(f *Forth) {
	f.DSPush(Cell(f.Base))
}

var dpaddr = func(f *Forth) {
	f.DSPush(Cell(f.DP))
}

var span = func(f *Forth) {
	f.DSPush(1000)
}

var source = func(f *Forth) {
	f.DSPush(Cell(len(f.TIB)))
	f.DSPush(0)
}

var lit = func(f *Forth) {
	val := f.Body[f.IP]
	f.IP++
	f.DSPush(val)
}

var exitword = func(f *Forth) {
	f.IP = len(f.Body)
}

var colon = func(f *Forth) {
	if f.State {
		forthError("ALREADY COMPILING")
	}
	token, ok := f.parseWord()
	if !ok {
		forthError("WORD NAME EXPECTED")
	}
	f.compileName = token
	f.State = true
	f.compileList = make([]Cell, 0, 64)
	f.compileWord = f.DefineWord(token, WordColon, nil, nil, false)
}

var semicolon = func(f *Forth) {
	if !f.State {
		forthError("NOT IN COMPILATION STATE")
	}
	exitXT, _ := f.FindXT("EXIT")
	f.compileList = append(f.compileList, Cell(exitXT))
	f.State = false
	bodyCopy := make([]Cell, len(f.compileList))
	copy(bodyCopy, f.compileList)
	if f.compileWord != nil {
		f.compileWord.Body = bodyCopy
	}
	f.compileList = f.compileList[:0]
	f.compileWord = nil
}

var lbracket = func(f *Forth) {
	f.State = false
}

var rbracket = func(f *Forth) {
	f.State = true
}

var state = func(f *Forth) {
	if f.State {
		f.DSPush(-1)
	} else {
		f.DSPush(0)
	}
}

var create = func(f *Forth) {
	token, ok := f.parseWord()
	if !ok {
		forthError("WORD NAME EXPECTED")
	}
	f.DefineWord(strings.ToUpper(token), WordCreate, nil, nil, false)
}

var constant = func(f *Forth) {
	val := f.DSPop()
	token, ok := f.parseWord()
	if !ok {
		forthError("WORD NAME EXPECTED")
	}
	f.DefineWord(strings.ToUpper(token), WordConstant, nil, []Cell{val}, false)
}

var variable = func(f *Forth) {
	token, ok := f.parseWord()
	if !ok {
		forthError("WORD NAME EXPECTED")
	}
	w := f.DefineWord(strings.ToUpper(token), WordVariable, nil, nil, false)
	f.Store(w.PFA, 0)
	f.DP = w.PFA + 2
}

var immediate = func(f *Forth) {
	w := f.LatestWord()
	if w != nil {
		w.Immediate = true
	}
}

var bracketcompile = func(f *Forth) {
	token, ok := f.parseWord()
	if !ok {
		forthError("WORD NAME EXPECTED")
	}
	w, found := f.Lookup(strings.ToUpper(token))
	if !found {
		f.EmitStr(token)
		f.EmitStr(" ?\r\n")
		return
	}
	_ = w
	xt, _ := f.FindXT(strings.ToUpper(token))
	f.compileList = append(f.compileList, Cell(xt))
}

var literal = func(f *Forth) {
	val := f.DSPop()
	litXT, _ := f.FindXT("LIT")
	f.compileList = append(f.compileList, Cell(litXT))
	f.compileList = append(f.compileList, val)
}

var tobody = func(f *Forth) {
	addr := int(f.DSPop())
	f.DSPush(Cell(addr + 2))
}

var parenword = func(f *Forth) {
	for f.IN < len(f.TIB) && f.TIB[f.IN] != ')' {
		f.IN++
	}
	if f.IN < len(f.TIB) {
		f.IN++
	}
}

var dotstring = func(f *Forth) {
	s := f.parseString()
	if f.State {
		addr := f.Here()
		for i := 0; i < len(s); i++ {
			f.CStore(addr+i, s[i])
		}
		f.DP = addr + len(s)
		litXT, _ := f.FindXT("LIT")
		f.compileList = append(f.compileList, Cell(litXT))
		f.compileList = append(f.compileList, Cell(addr))
		f.compileList = append(f.compileList, Cell(litXT))
		f.compileList = append(f.compileList, Cell(len(s)))
		typeXT, _ := f.FindXT("TYPE")
		f.compileList = append(f.compileList, Cell(typeXT))
	} else {
		f.EmitStr(s)
	}
}

var sstring = func(f *Forth) {
	s := f.parseString()
	if f.State {
		addr := f.Here()
		for i := 0; i < len(s); i++ {
			f.CStore(addr+i, s[i])
		}
		f.DP = addr + len(s)
		litXT, _ := f.FindXT("LIT")
		f.compileList = append(f.compileList, Cell(litXT))
		f.compileList = append(f.compileList, Cell(addr))
		f.compileList = append(f.compileList, Cell(litXT))
		f.compileList = append(f.compileList, Cell(len(s)))
	} else {
		addr := f.Here()
		for i := 0; i < len(s); i++ {
			f.CStore(addr+i, s[i])
		}
		f.DP = addr + len(s)
		f.DSPush(Cell(addr))
		f.DSPush(Cell(len(s)))
	}
}

var tick = func(f *Forth) {
	token, ok := f.parseWord()
	if !ok {
		forthError("WORD NAME EXPECTED")
	}
	w, found := f.Lookup(strings.ToUpper(token))
	if !found {
		forthError("WORD NOT FOUND: %s", token)
	}
	_ = w
	xt, _ := f.FindXT(strings.ToUpper(token))
	f.DSPush(Cell(xt))
}

var recurse = func(f *Forth) {
	if !f.State {
		forthError("RECURSE ONLY IN COMPILATION")
	}
	xt := len(f.Words) - 1
	f.compileList = append(f.compileList, Cell(xt))
}

var branch = func(f *Forth) {
	offset := int(f.Body[f.IP])
	f.IP += offset
}

var zerobranch = func(f *Forth) {
	flag := f.DSPop()
	offset := int(f.Body[f.IP])
	if flag == 0 {
		f.IP += offset
	} else {
		f.IP++
	}
}

var ifword = func(f *Forth) {
	if !f.State {
		forthError("IF ONLY IN COMPILATION")
	}
	xt, _ := f.FindXT("0BRANCH")
	f.compileList = append(f.compileList, Cell(xt))
	f.compileList = append(f.compileList, 0)
	f.DSPush(Cell(len(f.compileList) - 1))
}

var thenword = func(f *Forth) {
	if !f.State {
		forthError("THEN ONLY IN COMPILATION")
	}
	pos := int(f.DSPop())
	offset := Cell(len(f.compileList) - pos)
	f.compileList[pos] = offset
}

var elseword = func(f *Forth) {
	if !f.State {
		forthError("ELSE ONLY IN COMPILATION")
	}
	branchXT, _ := f.FindXT("BRANCH")
	f.compileList = append(f.compileList, Cell(branchXT))
	f.compileList = append(f.compileList, 0)
	pos2 := len(f.compileList) - 1
	pos1 := int(f.DSPop())
	offset1 := Cell(pos2 + 1 - pos1)
	f.compileList[pos1] = offset1
	f.DSPush(Cell(pos2))
}

var beginword = func(f *Forth) {
	if !f.State {
		forthError("BEGIN ONLY IN COMPILATION")
	}
	f.DSPush(Cell(len(f.compileList)))
}

var untilword = func(f *Forth) {
	if !f.State {
		forthError("UNTIL ONLY IN COMPILATION")
	}
	dest := int(f.DSPop())
	xt, _ := f.FindXT("0BRANCH")
	f.compileList = append(f.compileList, Cell(xt))
	offset := Cell(dest - len(f.compileList))
	f.compileList = append(f.compileList, offset)
}

var againword = func(f *Forth) {
	if !f.State {
		forthError("AGAIN ONLY IN COMPILATION")
	}
	dest := int(f.DSPop())
	xt, _ := f.FindXT("BRANCH")
	f.compileList = append(f.compileList, Cell(xt))
	offset := Cell(dest - len(f.compileList))
	f.compileList = append(f.compileList, offset)
}

var whileword = func(f *Forth) {
	if !f.State {
		forthError("WHILE ONLY IN COMPILATION")
	}
	dest := int(f.DSPop())
	xt, _ := f.FindXT("0BRANCH")
	f.compileList = append(f.compileList, Cell(xt))
	f.compileList = append(f.compileList, 0)
	f.DSPush(Cell(dest))
	f.DSPush(Cell(len(f.compileList) - 1))
}

var repeatword = func(f *Forth) {
	if !f.State {
		forthError("REPEAT ONLY IN COMPILATION")
	}
	orig := int(f.DSPop())
	dest := int(f.DSPop())
	branchXT, _ := f.FindXT("BRANCH")
	f.compileList = append(f.compileList, Cell(branchXT))
	back := Cell(dest - len(f.compileList))
	f.compileList = append(f.compileList, back)
	forward := Cell(len(f.compileList) - orig)
	f.compileList[orig] = forward
}
