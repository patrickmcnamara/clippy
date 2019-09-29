package clippy

import (
	"fmt"
	"strings"
	"unicode"
)

// EmptyValue is an empty value. This is used for flags where the default value should be the empty string.
var EmptyValue = "\000"

// Flag is a string value given in the command line (or by a default value).
type Flag struct {
	Name         string // Name of the flag.
	Alias        rune   // Alias of the flag.
	Type         string // Type of the flag. For example, "FILENAME" or "URL".
	Description  string // Description of the flag.
	DefaultValue string // Default value of the flag. If it is left empty, it is assumed that the flag is mandatory and must be given by the user. Use EmptyValue if the default value should be empty.
}

func (f *Flag) String() string {
	var sb strings.Builder
	sb.WriteString("--" + f.Name)
	if f.Alias != rune(0) {
		sb.WriteString(", -" + string(f.Alias))
	}
	sb.WriteString("\t" + f.Description)
	if f.DefaultValue == "_" {
		sb.WriteString(fmt.Sprintf(" (%q)", ""))
	} else if f.DefaultValue != "" {
		sb.WriteString(fmt.Sprintf(" (%q)", f.DefaultValue))
	}
	return sb.String()
}

func (f *Flag) check() error {
	// Check that the flag has a name.
	if f.Name == "" {
		return fmt.Errorf("missing name of flag")
	}

	// Check that each character in the flag's name is valid.
	for _, char := range f.Name {
		if !unicode.IsLetter(char) && !unicode.IsNumber(char) && char != '-' {
			return fmt.Errorf("invalid character in flag name: %q", char)
		}
	}

	// Check if the flag's alias is a valid.
	if f.Alias != rune(0) {
		if !unicode.IsLetter(f.Alias) && !unicode.IsNumber(f.Alias) {
			return fmt.Errorf("flag alias is an invalid character: %q", f.Alias)
		}
	}

	return nil
}

// FlagSet is a list of Flags.
type FlagSet []*Flag

func (fs *FlagSet) String() string {
	var sb strings.Builder
	for _, f := range *fs {
		sb.WriteString("\t" + f.String())
	}
	return sb.String()
}

func (fs *FlagSet) check() error {
	names := make(map[string]struct{})
	for _, f := range *fs {
		// Check the flag.
		if err := f.check(); err != nil {
			return err
		}

		// Check if the flag's name already exists.
		if _, ok := names[f.Name]; !ok {
			names[f.Name] = struct{}{}
		} else {
			return fmt.Errorf("duplicate flag name or alias: %q", f.Name)
		}

		// Check if the flag's alias already exists.
		if f.Alias != rune(0) {
			if _, ok := names[string(f.Alias)]; !ok {
				names[string(f.Alias)] = struct{}{}
			} else {
				return fmt.Errorf("duplicate flag name or alias: %q", f.Alias)
			}
		}
	}
	return nil
}

func (fs *FlagSet) get(name string) *Flag {
	for _, flag := range *fs {
		if "--"+flag.Name == name || flag.Alias != rune(0) && "-"+string(flag.Alias) == name {
			return flag
		}
	}
	return nil
}

func (fs *FlagSet) parse(params []string) (flags map[string]string, args []string, err error) {
	flags = make(map[string]string)
	args = make([]string, 0)

	// Parse given flag values and arguments.
	for i := 0; i < len(params); i++ {
		param := params[i]
		if flag := fs.get(param); flag != nil {
			if i+1 < len(params) {
				flags[flag.Name] = params[i+1]
				i++
			} else {
				err = fmt.Errorf("no corresponding value for flag: %q", param)
				return
			}
		} else {
			args = append(args, param)
		}
	}

	// Check for default flag values.
	for _, f := range *fs {
		name := f.Name
		if _, ok := flags[f.Name]; !ok {
			if f.DefaultValue == "" {
				err = fmt.Errorf("no given or default value for flag: %q", name)
				return
			} else if f.DefaultValue == EmptyValue {
				flags[name] = ""
			} else {
				flags[name] = f.DefaultValue
			}
		}
	}

	return
}
