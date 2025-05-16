# pforth

FORTH interpreter and compiler implemented in Go, ported from the original Z80 and MOS 6502 assembly implementations.

## Highlights

- Forth interpreter/compiler with interactive REPL
- Colon definitions with full compiler support (IF/THEN/ELSE, BEGIN/UNTIL/WHILE/REPEAT, DO/LOOP)
- Byte-addressable 64KB memory model with separate data and return stacks
- Recursive definitions via RECURSE
- String output with ."
- 80+ primitive words covering stack, memory, arithmetic, logic, I/O, and compiler operations
- Dual-layer dictionary: Go primitives plus Forth words loaded from `libs/core.fs`
- Demo programs for sorting, factorial, Fibonacci, sieve, and Tower of Hanoi

## Prerequisites

- **Go 1.23+** — required to build from source; [download](https://go.dev/dl/)

## Installation

### Build from Source

```bash
git clone https://github.com/carlosrabelo/pforth.git
cd pforth-go
make build
```

Install to `~/.local/bin` (default), or system-wide to `/usr/local/bin` (sudo only for the copy):

```bash
make install
make install-system
make uninstall
make uninstall-system
```

## Usage

```bash
make run
```

Or directly:

```bash
./bin/pforth
```

Load a program, then enter the REPL:

```bash
./bin/pforth demos/hello.fs
```

### Example session

```forth
1 2 + .
3 ok

: SQUARE DUP * ;
ok
5 SQUARE .
25 ok

: FACTORIAL
  DUP 0= IF DROP 1 ELSE DUP 1- RECURSE * THEN ;
ok
10 FACTORIAL .
3628800 ok
```

## Documentation

- [Architecture](docs/architecture.md) — inner interpreter, memory model, execution engine
- [Word Reference](docs/words.md) — complete reference of all built-in words

## Project Layout

```
pforth/              # All Go source code
├── cmd/pforth/      # Entry point
└── internal/forth/  # Core engine: stacks, memory, dictionary, primitives
libs/                # Forth libraries (core.fs)
demos/               # Sample Forth programs
docs/                # Architecture and word reference
bin/                 # Compiled binaries (git-ignored)
.make/               # Build and install scripts
```

## Development

```bash
make build             # Compile binary to bin/pforth
make run               # Build and run the REPL
make test              # Run all tests
make quality           # Format, vet, and lint
make install           # Install binary to ~/.local/bin
make install-system    # Install binary to /usr/local/bin
make uninstall         # Remove from ~/.local/bin
make uninstall-system  # Remove from /usr/local/bin
```

## License

This project is licensed under the MIT License — see [LICENSE](LICENSE) for details.
