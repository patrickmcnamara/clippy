package clippy

import (
	"fmt"
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
			}
		}

		// Action, parse and setup error handler using chkErr.
		actionErr = func(err error) { chkErr(err, ActionErrHandler) }
		parseErr  = func(err error) { chkErr(err, ParseErrHandler) }
		setupErr  = func(err error) { chkErr(err, SetupErrHandler) }
	)

	// Check for errors with commands and flags.
	setupErr(c.Check())

	// Run subcommand or help or version if it's there.
	if len(params) >= 1 {
		p1 := params[0]
		if command := c.Commands.get(p1); command != nil {
			parseErr(command.run(c.Name, params[1:]))
			return
		} else if p1 == "-h" || p1 == "--help" {
			fmt.Println(c.help())
			return
		} else if p1 == "-v" || p1 == "--version" {
			fmt.Println(c.version())
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
	usage := "[global flags...] [command] [flags and values...] [arguments...]"
	if c.Usage != "" {
		usage = c.Usage
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
		sb.WriteString(c.Commands.help("\t"))
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
		sb.WriteString(c.Flags.help("\t"))
		sb.WriteRune('\n')
	}

	return strings.TrimRight(sb.String(), "\n")
}
