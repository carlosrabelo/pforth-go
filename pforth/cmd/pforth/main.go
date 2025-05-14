package main

import (
	"os"

	"github.com/carlosrabelo/pforth/pforth/internal/forth"
)

func main() {
	f := forth.New(os.Stdin, os.Stdout)

	for _, arg := range os.Args[1:] {
		f.LoadFile(arg)
	}

	f.QUIT()
}
