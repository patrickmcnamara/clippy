package clippy

import (
	"fmt"
	"strings"
	"unicode"
)

// Command is a subcommand for a program.
type Command struct {
	Names       []string // Name and aliases of the command. It is required.
	Description string   // Description of the command.
	Usage       string   // Usage describes how to use the command. It has a default.
	Flags       FlagSet  // Flags used by the program.
	Action      Action   // Action is called when this particular command is.
}

func (c *Command) check() error {
	// Check that the command has at least one name.
	if len(c.Names) < 1 {
		return fmt.Errorf("missing name of command")
	}

	// Check that each character in each name is valid.
	for _, name := range c.Names {
		for _, char := range name {
			if !unicode.IsLetter(char) && !unicode.IsNumber(char) && char != '-' {
				return fmt.Errorf("invalid character in command name: %q in %q", char, name)
			}
		}
	}

	// Check each flag in command's flagset.
	for _, f := range c.Flags {
		if err := f.check(); err != nil {
			return err
		}
	}

	return nil
}

func (c *Command) run(name string, params []string) error {
	// Check for help flag.
	if len(params) >= 1 && params[0] == "-h" || params[0] == "--help" {
		fmt.Println(c.help(name))
		return nil
	}

	// Parse parameters for flags and arguments.
	flags, args, err := c.Flags.parse(params)
	if err != nil {
		return err
	}

	// Check if there is a default action.
	if c.Action == nil {
		return DefaultAction(flags, args)
	}

	// Run action if there is one.
	return c.Action(flags, args)
}

func (c *Command) help(name string) string {
	var sb strings.Builder

	// NAME
	sb.WriteString("NAME:\n")
	sb.WriteString("\t" + name + " " + c.Names[0])
	sb.WriteString("\n\n")

	// DESCRIPTION
	if c.Description != "" {
		sb.WriteString("DESCRIPTION:\n")
		sb.WriteString("\t" + c.Description + "\n\n")
	}

	// USAGE
	sb.WriteString("USAGE:\n")
	usage := "[flags and values...] [arguments...]"
	if c.Usage != "" {
		usage = c.Usage
	}
	sb.WriteString("\t" + name + " " + c.Names[0] + " " + usage + "\n\n")

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

// CommandSet is a list of Commands.
type CommandSet []*Command

func (cs *CommandSet) check() error {
	names := make(map[string]struct{})
	for _, command := range *cs {
		// Check the command.
		if err := command.check(); err != nil {
			return err
		}

		// Check that each of the command's names is unique.
		for _, name := range command.Names {
			if _, ok := names[name]; !ok {
				names[name] = struct{}{}
			} else {
				return fmt.Errorf("duplicate command name %q", name)
			}
		}
	}
	return nil
}

func (cs *CommandSet) get(name string) *Command {
	for _, command := range *cs {
		for _, commandName := range command.Names {
			if commandName == name {
				return command
			}
		}
	}
	return nil
}

func (cs *CommandSet) help(indent string) string {
	var sb strings.Builder

	var width int
	var names []string
	for _, cmd := range *cs {
		name := strings.Join(cmd.Names, ", ")
		if l := len(name); l > width {
			width = l
		}
		names = append(names, name)
	}

	for i, cmd := range *cs {
		sb.WriteString(fmt.Sprintf("%s%-*s%s%s\n", indent, width, names[i], indent, cmd.Description))
	}

	return sb.String()
}
