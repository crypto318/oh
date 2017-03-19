// Released under an MIT license. See LICENSE.

package cell

import "errors"

type ApplyModer interface {
	ApplyMode() error
}

type Engine interface {
	Deref(name string, address uintptr) Cell
	MakeParser(ReadStringer, func(Cell, string, int) (Cell, bool)) Parser
	Throw(filename string, lineno int, message string)
}

type Parser interface {
	Interpret(string) bool
	ParsePipe() (bool, interface{})
	State(string) (string, string, string)
}

type ParserMaker func(
	Engine, ReadStringer,
	func(Cell, string, int) (Cell, bool),
) Parser

type ReadStringer interface {
	ReadString(delim byte) (line string, err error)
}

type UI interface {
	Close() error
	Exists() bool
	ReadString(delim byte) (string, error)
	TerminalMode() (ApplyModer, error)
}

type UIMaker func([]string) UI

const (
	ErrNotExecutable = "oh: 126: error/runtime: "
	ErrNotFound      = "oh: 127: error/runtime: "
	ErrSyntax        = "oh: 1: error/syntax: "
)

var ErrCtrlCPressed = errors.New("ctrl-c pressed")
