# Architecture

Overview of the pforth inner interpreter, memory model, and execution engine.

## Stacks

Two separate LIFO stacks, implemented as Go slices of `Cell` (int):

| Stack | Variable | Purpose |
|---|---|---|
| Data Stack (DS) | `f.DS` | Operand storage for arithmetic, memory, and logic words |
| Return Stack (RS) | `f.RS` | Loop indices, colon definition return addresses, `>R`/`R>` temporary storage |

Stack operations panic with a descriptive message on underflow.

## Memory Model

Flat byte array of 65536 elements (`[]byte`):

- **Cells** are 16-bit values stored in little-endian format (two bytes per cell).
- **Dictionary Pointer (DP)** starts at 1024, reserving the first 1KB for system use.
- **`@` (fetch)** reads 2 bytes and returns a `Cell`.
- **`!` (store)** writes a `Cell` as 2 bytes.
- **`C@` / `C!`** read/write single bytes.
- **`HERE`** returns the current DP value.
- **`ALLOT`** advances DP by n bytes.
- **`,` (comma)** stores a cell at DP and advances by 2.
- **`C,` (c-comma)** stores a byte at DP and advances by 1.

Byte/word access is bounds-checked and panics on out-of-range addresses.

## Dictionary

The dictionary is a slice of `*Word` pointers, plus a `map[string]int` for name-to-index lookup:

```go
type Word struct {
    Name      string
    Immediate bool
    Type      WordType    // WordPrimitive, WordColon, WordConstant, WordVariable, WordCreate
    Code      XTCode      // Native Go function (primitives only)
    Body      []Cell      // Execution tokens for colon definitions
    PFA       int         // Parameter field address (VARIABLE/CREATE)
}
```

### Word Types

| Type | Behavior |
|---|---|
| `WordPrimitive` | Calls `w.Code(f)` — a native Go function |
| `WordColon` | Saves `f.Body`/`f.IP`, runs `w.Body[]` cell by cell, restores state |
| `WordConstant` | Pushes `w.Body[0]` onto DS |
| `WordVariable` | Pushes `w.PFA` onto DS |
| `WordCreate` | Pushes `w.PFA` onto DS |

### Name Resolution

- All names are normalized to uppercase via `strings.ToUpper` during definition and lookup.
- `FIND` returns the XT (index into the `f.Words` slice) and a success flag.
- `'` (tick) returns the XT of the next parsed word.
- `>BODY` converts an XT to its parameter field address (XT + 2 cell offsets).

## Execution Engine

### Colon Definitions

When `ExecuteWord` encounters a `WordColon`:

1. Save current `f.Body` and `f.IP` on a Go call stack (local variables).
2. Set `f.Body = w.Body` and `f.IP = 0`.
3. Loop: read `f.Body[f.IP]`, increment IP, call `ExecuteWord` on the token.
4. Restore previous `f.Body`/`f.IP`.

This approach is recursive in Go — each nested colon definition adds a frame. The same mechanism supports `EXIT` (set IP past end → loop exits, previous frame restored).

### Branching

Branch offsets are stored inline in the body after the branch XT:

- **`BRANCH`**: `f.IP += f.Body[f.IP]` (unconditional jump)
- **`0BRANCH`**: if TOS = 0, `f.IP += f.Body[f.IP]`; else `f.IP++` (skip offset)

Compile-time words (`IF`, `THEN`, `ELSE`, `BEGIN`, `UNTIL`, `AGAIN`, `WHILE`, `REPEAT`) patch forward/backward offsets during compilation using the data stack as a temporary address tracker.

### DO/LOOP

Runtime state stored on the return stack as pairs: `( limit index -- )`.

- **`(DO)`**: pops `limit start` from DS, pushes `limit start` onto RS.
- **`I`**: copies the top index from RS to DS.
- **`J`**: reads the index from the next RS pair down (nested loops).
- **`(LOOP)`**: increments index, if `index < limit` jumps back by the inline offset.
- **`(+LOOP)`**: adds a step value to index, tests for loop termination.
- **`UNLOOP`**: removes one RS pair (used by `LEAVE`).

## Compilation State

When `STATE` is true (`f.State`), tokens are appended to `f.compileList` as XTs (execution tokens) instead of being executed immediately. Numbers are compiled as `LIT <value>` pairs.

- **`:`** — creates a new `WordColon` placeholder, sets `STATE` to true, records name.
- **`;`** — appends `EXIT`, finalizes the body, sets `STATE` to false, resets compile list.
- **`IMMEDIATE`** — marks the most recently defined word as immediate (executed even in compile state).

## Interpreter Loop

```
QUIT → (
  if interpreting: print "OK"
  read a line into TIB
  set IN = 0
  interpretLoop()
    parse a token
    look up in dictionary
    if found:
      if compiling AND not immediate: compile XT
      else: execute word
    else: try to parse as number
    if failed: print " ?", abort
)
```

The text input buffer (`TIB`) and input pointer (`IN`) support re-parsing with `WORD` and `EXPECT`.

## REPL

The REPL is driven by `QUIT`, which loops until `f.Running` becomes false (set by `BYE` or EOF on stdin). Each iteration reads a line via `bufio.Reader.ReadString('\n')`, so there is no fixed line-length limit.
