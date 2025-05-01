package forth

import "fmt"

type Cell int

type WordType int

const (
	WordPrimitive WordType = iota
	WordColon
	WordConstant
	WordVariable
	WordCreate
)

type XTCode func(f *Forth)

type Word struct {
	Name      string
	Immediate bool
	Type      WordType
	Code      XTCode
	Body      []Cell
	PFA       int
}

type ForthError struct {
	Message string
}

func (e *ForthError) Error() string {
	return e.Message
}

func forthError(format string, args ...any) {
	panic(&ForthError{Message: fmt.Sprintf(format, args...)})
}
