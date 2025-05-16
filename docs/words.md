# Word Reference

Complete reference of all built-in pforth words, organized by category. Stack notation follows the Forth convention `( before -- after )` where top-of-stack is rightmost.

## Stack Manipulation

| Word | Stack Effect | Description |
|---|---|---|
| `DUP` | `( n -- n n )` | Duplicate top of stack |
| `?DUP` | `( n -- 0 \| n n )` | Duplicate only if non-zero |
| `DROP` | `( n -- )` | Remove top of stack |
| `SWAP` | `( a b -- b a )` | Exchange top two items |
| `OVER` | `( a b -- a b a )` | Copy second item to top |
| `ROT` | `( a b c -- b c a )` | Rotate third item to top |
| `-ROT` | `( a b c -- c a b )` | Rotate top item to third |
| `NIP` | `( a b -- b )` | Drop second item |
| `TUCK` | `( a b -- b a b )` | Duplicate top under second |
| `2DUP` | `( a b -- a b a b )` | Duplicate top pair |
| `2DROP` | `( a b -- )` | Drop top pair |
| `2SWAP` | `( a b c d -- c d a b )` | Exchange top two pairs |
| `2OVER` | `( a b c d -- a b c d a b )` | Copy third pair to top |
| `>R` | `( n -- )` ( R: `-- n` ) | Move from data to return stack |
| `R>` | `( -- n )` ( R: `n --` ) | Move from return to data stack |
| `R@` | `( -- n )` ( R: `n -- n` ) | Copy return stack top to data stack |
| `DEPTH` | `( -- +n )` | Number of items on data stack |

## Memory

| Word | Stack Effect | Description |
|---|---|---|
| `@` | `( addr -- n )` | Fetch 16-bit cell from memory |
| `!` | `( n addr -- )` | Store 16-bit cell to memory |
| `C@` | `( addr -- c )` | Fetch byte from memory |
| `C!` | `( c addr -- )` | Store byte to memory |
| `+!` | `( n addr -- )` | Add n to memory cell |
| `?` | `( addr -- )` | Fetch and print value at address |
| `FILL` | `( addr len c -- )` | Fill memory region with byte |
| `ERASE` | `( addr len -- )` | Fill memory region with zeros |
| `CMOVE` | `( dst src len -- )` | Copy bytes forward |
| `CMOVE>` | `( dst src len -- )` | Copy bytes backward |

## Arithmetic

| Word | Stack Effect | Description |
|---|---|---|
| `+` | `( a b -- sum )` | Add |
| `-` | `( a b -- diff )` | Subtract (b - a) |
| `*` | `( a b -- prod )` | Multiply |
| `/` | `( a b -- quot )` | Divide (b / a) |
| `MOD` | `( a b -- rem )` | Remainder (b mod a) |
| `/MOD` | `( a b -- rem quot )` | Divide, return remainder and quotient |
| `NEGATE` | `( n -- -n )` | Two's complement negation |
| `ABS` | `( n -- |n| )` | Absolute value |
| `MIN` | `( a b -- min )` | Minimum |
| `MAX` | `( a b -- max )` | Maximum |
| `1+` | `( n -- n+1 )` | Increment by 1 |
| `1-` | `( n -- n-1 )` | Decrement by 1 |
| `2+` | `( n -- n+2 )` | Increment by 2 |
| `2-` | `( n -- n-2 )` | Decrement by 2 |
| `2*` | `( n -- n*2 )` | Multiply by 2 (left shift) |
| `2/` | `( n -- n/2 )` | Divide by 2 (right shift) |

## Comparison and Logic

| Word | Stack Effect | Description |
|---|---|---|
| `=` | `( a b -- flag )` | Equal |
| `<>` | `( a b -- flag )` | Not equal |
| `<` | `( a b -- flag )` | Signed less than |
| `>` | `( a b -- flag )` | Signed greater than |
| `<=` | `( a b -- flag )` | Signed less than or equal |
| `>=` | `( a b -- flag )` | Signed greater than or equal |
| `0=` | `( n -- flag )` | Zero test |
| `0<>` | `( n -- flag )` | Non-zero test |
| `0<` | `( n -- flag )` | Negative test |
| `0>` | `( n -- flag )` | Positive test |
| `U<` | `( a b -- flag )` | Unsigned less than |
| `U>` | `( a b -- flag )` | Unsigned greater than |
| `U<=` | `( a b -- flag )` | Unsigned less than or equal |
| `U>=` | `( a b -- flag )` | Unsigned greater than or equal |
| `AND` | `( a b -- c )` | Bitwise AND |
| `OR` | `( a b -- c )` | Bitwise OR |
| `XOR` | `( a b -- c )` | Bitwise XOR |
| `INVERT` | `( n -- ~n )` | Bitwise NOT |

Flags return -1 (true) or 0 (false), following Forth convention.

## I/O

| Word | Stack Effect | Description |
|---|---|---|
| `KEY` | `( -- c )` | Read one character from stdin |
| `EMIT` | `( c -- )` | Write one character to stdout |
| `CR` | `( -- )` | Write carriage return + newline |
| `SPACE` | `( -- )` | Write a space |
| `SPACES` | `( n -- )` | Write n spaces |
| `TYPE` | `( addr len -- )` | Write len bytes from addr |
| `.` | `( n -- )` | Print signed integer in base |
| `U.` | `( u -- )` | Print unsigned integer in base |
| `."` | `( -- )` (immediate) | Print string literal during interpretation |
| `PAGE` | `( -- )` | Clear screen (emits form feed) |
| `BYE` | `( -- )` | Exit pforth |

## Interpreter / Misc

| Word | Stack Effect | Description |
|---|---|---|
| `WORD` | `( delim -- addr )` | Parse next word delimited by char, return counted string address |
| `FIND` | `( addr -- xt flag )` | Look up counted string in dictionary |
| `EXECUTE` | `( xt -- )` | Execute word by execution token |
| `INTERPRET` | `( -- )` | Interpret text in TIB from current IN position |
| `QUIT` | `( -- )` | Reset stacks and enter REPL loop |
| `ABORT` | `( -- )` | Reset stacks, stop compilation, discard TIB |
| `EXPECT` | `( addr len -- )` | Read line from stdin into buffer, stores count at address 1000 |
| `.S` | `( -- )` | Print data stack contents with depth markers `< ... >` |
| `STATE` | `( -- flag )` | Current interpreter state (-1 compiling, 0 interpreting) |
| `BASE` | `( -- addr )` | Address of base variable |
| `DP` | `( -- addr )` | Address of dictionary pointer variable |
| `SPAN` | `( -- addr )` | Address of span variable (stores EXPECT char count) |
| `HERE` | `( -- addr )` | Current dictionary pointer |
| `PAD` | `( -- addr )` | Temporary buffer address (HERE + 68) |

## Compiler Words

### Defining Words

| Word | Stack Effect | Description |
|---|---|---|
| `:` | `( -- )` (immediate) | Start colon definition |
| `;` | `( -- )` (immediate) | End colon definition |
| `CONSTANT` | `( n -- )` | Create named constant with value n |
| `VARIABLE` | `( -- )` | Create named variable (2 bytes reserved) |
| `CREATE` | `( -- )` | Create named dictionary header, no size allocated |
| `IMMEDIATE` | `( -- )` | Mark latest word as immediate |

### Compiler Directives

| Word | Stack Effect | Description |
|---|---|---|
| `[` | `( -- )` (immediate) | Enter interpretation state |
| `]` | `( -- )` (immediate) | Enter compilation state |
| `LITERAL` | `( n -- )` (immediate) | Compile the top-of-stack as a literal |
| `[COMPILE]` | `( -- )` (immediate) | Force compilation of next immediate word |
| `'` | `( -- xt )` | Tick — get XT of next word |
| `>BODY` | `( xt -- pfa )` | Convert XT to parameter field address |
| `RECURSE` | `( -- )` (immediate) | Recursively compile current definition |
| `S"` | `( -- addr len )` (immediate) | String literal — pushes addr len at runtime |

### Memory Allocation

| Word | Stack Effect | Description |
|---|---|---|
| `ALLOT` | `( n -- )` | Advance DP by n bytes |
| `,` | `( n -- )` | Compile one cell at HERE, advance DP by 2 |
| `C,` | `( c -- )` | Compile one byte at HERE, advance DP by 1 |

## Branching

### Conditional

| Word | Stack Effect | Description |
|---|---|---|
| `IF` | `( -- )` (immediate) | Begin conditional branch |
| `THEN` | `( -- )` (immediate) | End conditional branch |
| `ELSE` | `( -- )` (immediate) | Alternative branch |

```
IF <true-part> THEN
IF <true-part> ELSE <false-part> THEN
```

### Indefinite Loops

| Word | Stack Effect | Description |
|---|---|---|
| `BEGIN` | `( -- )` (immediate) | Start loop |
| `UNTIL` | `( -- )` (immediate) | Loop until flag true |
| `AGAIN` | `( -- )` (immediate) | Loop unconditionally |
| `WHILE` | `( -- )` (immediate) | Conditional exit of BEGIN loop |
| `REPEAT` | `( -- )` (immediate) | End BEGIN/WHILE/REPEAT loop |

```
BEGIN ... UNTIL
BEGIN ... AGAIN
BEGIN ... WHILE ... REPEAT
```

### Definite Loops

| Word | Stack Effect | Description |
|---|---|---|
| `DO` | `( limit start -- )` (immediate) | Start counted loop |
| `LOOP` | `( -- )` (immediate) | End counted loop, step = 1 |
| `+LOOP` | `( step -- )` (immediate) | End counted loop, step = n |
| `I` | `( -- index )` | Current loop index |
| `J` | `( -- index )` | Outer loop index |
| `LEAVE` | `( -- )` (immediate) | Exit loop early |
| `UNLOOP` | `( -- )` | Discard loop parameters from return stack |

```
10 0 DO I . LOOP         \ prints 0 1 2 ... 9
10 0 DO I . 3 +LOOP      \ prints 0 3 6 9
```

## Runtime Primitives

These words are not typically called directly; they are emitted by the compiler:

| Word | Description |
|---|---|
| `LIT` | Push next inline cell as a literal |
| `EXIT` | Return from colon definition |
| `BRANCH` | Unconditional relative jump |
| `0BRANCH` | Conditional relative jump (zero) |
| `(DO)` | Initialize DO/LOOP parameters on return stack |
| `(LOOP)` | Increment index and loop if not done |
| `(+LOOP)` | Add step and loop if not done |
