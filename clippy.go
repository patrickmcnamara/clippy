package clippy

import (
	"fmt"
	"os"
	"strings"
)

// Clippy represents a CLI program.
type Clippy struct {
	Name        string     // Name of the program. It is required.
	Tagline     string     // Tagline of the program.
	Version     string     // Version of the program. It is required.
	Description string     // Description of the program.
	Authors     []Author   // A list of authors of the program.
	Usage       string     // Usage describes how to use the program. It has a default.
	Flags       FlagSet    // Global flags used by the program.
	Commands    CommandSet // Commands are the subcommands of the program.
	Action      Action     // Action is called when this particular command is.
}

// Run checks the clippy setup, parses params and runs the parsed command, handling errors it encounters.
func (c *Clippy) Run(params []string) {
	// Establish error handlers.
	var (
		// Only use error handler if err != nil.
		chkErr = func(err error, errHandler ErrHandler) {
			if err != nil {
				errHandler(c.Name, err)
				os.Exit(4)
			}
		}

		// Action, parse and setup error handler using chkErr.
		actionErr = func(err error) { chkErr(err, ActionErrHandler) }
		parseErr  = func(err error) { chkErr(err, ParseErrHandler) }
		setupErr  = func(err error) { chkErr(err, SetupErrHandler) }
	)

	// Check for errors with commands and flags.
	setupErr(c.Check())

	// Check for global flags.
	for _, parameter := range params {
		switch parameter {
		case "--help", "-h":
			fmt.Println(c.help())
			return
		case "--version", "-v":
			fmt.Println(c.version())
			return
		}
	}

	// Run subcommand if it's there.
	if len(params) >= 1 {
		if command := c.Commands.get(params[0]); command != nil {
			parseErr(command.run(params[1:]))
			return
		}
	}

	// Parse flags and arguments.
	flags, args, err := c.Flags.parse(params)
	parseErr(err)

	// Run default action if none is set.
	if c.Action == nil {
		parseErr(HelpAction(flags, args))
	}
	// Otherwise run given action.
	actionErr(c.Action(flags, args))
}

// Check checks clippy.
func (c *Clippy) Check() error {
	// Check for errors with flags.
	if err := c.Flags.check(); err != nil {
		return err
	}

	// Check for errors with commands.
	if err := c.Commands.check(); err != nil {
		return err
	}

	return nil
}

func (c *Clippy) String() string {
	return c.version()
}

func (c *Clippy) version() string {
	return c.Name + " " + c.Version
}

func (c *Clippy) help() string {
	var sb strings.Builder

	// NAME and TAGLINE
	sb.WriteString("NAME:\n")
	sb.WriteString("\t" + c.Name)
	if c.Tagline != "" {
		sb.WriteString(" - " + c.Tagline)
	}
	sb.WriteString("\n\n")

	// VERSION
	sb.WriteString("VERSION:\n")
	sb.WriteString("\t" + c.Version + "\n\n")

	// DESCRIPTION
	if c.Description != "" {
		sb.WriteString("DESCRIPTION:\n")
		sb.WriteString("\t" + c.Description + "\n\n")
	}

	// AUTHOR(S)
	if len(c.Authors) >= 1 {
		sb.WriteString("AUTHOR")
		if len(c.Authors) > 1 {
			sb.WriteString("S:\n")
		} else {
			sb.WriteString(":\n")
		}
		for _, author := range c.Authors {
			sb.WriteString("\t" + author.String() + "\n")
		}
		sb.WriteRune('\n')
	}

	// USAGE
	sb.WriteString("USAGE:\n")
	var usage string
	if c.Usage != "" {
		usage = c.Usage
	} else {
		usage = "[global flags...] [command] [flags and values...] [arguments...]"
	}
	sb.WriteString("\t" + c.Name + " " + usage + "\n\n")

	// GLOBAL FLAGS
	sb.WriteString("GLOBAL FLAGS:\n")
	sb.WriteString("\t" + "--help, -h  \tshow help (with optional subcommand) and exit\n")
	sb.WriteString("\t" + "--version, -v  \tshow version and exit\n")
	sb.WriteRune('\n')

	// COMMANDS
	if len(c.Commands) >= 1 {
		sb.WriteString("COMMAND")
		if len(c.Commands) > 1 {
			sb.WriteString("S:\n")
		} else {
			sb.WriteString(":\n")
		}
		for _, command := range c.Commands {
			sb.WriteString("\t" + command.String() + "\n")
		}
		sb.WriteRune('\n')
	}

	// FLAGS
	if len(c.Flags) >= 1 {
		sb.WriteString("FLAG")
		if len(c.Flags) > 1 {
			sb.WriteString("S:\n")
		} else {
			sb.WriteString(":\n")
		}
		for _, flag := range c.Flags {
			sb.WriteString("\t" + flag.String() + "\n")
		}
		sb.WriteRune('\n')
	}

	return strings.TrimRight(sb.String(), "\n")
}
