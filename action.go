package clippy

import "errors"

// Action represents a function run by a command.
type Action func(flags map[string]string, args []string) error

// DefaultAction is a no-op. It does nothing at all.
var DefaultAction Action = func(flags map[string]string, args []string) error { return nil }

// HelpAction returns an error saying to use the "--help" global flag.
var HelpAction Action = func(flags map[string]string, args []string) error {
	return errors.New("use the \"--help\" global flag")
}
