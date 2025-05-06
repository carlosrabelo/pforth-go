package forth

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
