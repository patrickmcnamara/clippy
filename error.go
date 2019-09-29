package clippy

import (
	"fmt"
	"os"
)

// ErrHandler is an error handler that handles an action, parsing, or setup error.
// The name is the name of the program and the err is the error being handled.
type ErrHandler func(name string, err error)

var (
	// ActionErrHandler handles errors the errors that actions may return.
	ActionErrHandler ErrHandler = func(name string, err error) { defaultErrHandler(name, err, 1) }
	// ParseErrHandler handles errors the errors that may be returned when parsing.
	ParseErrHandler ErrHandler = func(name string, err error) { defaultErrHandler(name, err, 2) }
	// SetupErrHandler handles errors the errors that may be returned when checking if the clippy is valid.
	SetupErrHandler ErrHandler = func(name string, err error) { defaultErrHandler(name, err, 3) }
)

func defaultErrHandler(name string, err error, exitCode int) {
	fmt.Fprintln(os.Stderr, fmt.Errorf("%s: %w", name, err))
	os.Exit(exitCode)
}
