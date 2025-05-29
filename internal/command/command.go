package command

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"unicode"
)

type Command struct {
	command string
	args    []string
	opts    []rune
}

func (c *Command) String() string {
	return fmt.Sprintf("cmd : %s, args : %v, opts : %s", c.command, c.args, string(c.opts))
}

func (c *Command) GetCommand() string {
	return c.command
}

func (c *Command) GetArgs() []string {
	return c.args
}

func (c *Command) GetOpts() []rune {
	return c.opts
}

func (c *Command) HasOpt(opt rune) bool {
	for _, o := range c.opts {
		if o == opt {
			return true
		}
	}
	return false
}

func ParseCommand(command string) *Command {
	// Handle empty commands
	command = strings.TrimSpace(command)
	if command == "" {
		return &Command{}
	}

	commandList := parseCommandLine(command)

	if len(commandList) == 0 {
		return &Command{}
	}

	cmd := Command{
		command: commandList[0],
		args:    commandList[1:],
	}

	var filteredArgs []string
	for _, arg := range cmd.args {
		if len(arg) > 0 && arg[0] == '-' {
			// Handle both short (-abc) and long (--flag) options
			if len(arg) > 1 && arg[1] == '-' {
				// Long option like --verbose, store as single option
				cmd.opts = append(cmd.opts, []rune(arg[2:])...)
			} else {
				// Short options like -abc, each letter is an option
				cmd.opts = append(cmd.opts, []rune(arg[1:])...)
			}
		}
		filteredArgs = append(filteredArgs, arg)
	}
	cmd.args = filteredArgs

	return &cmd
}

// parseCommandLine splits a command line into arguments, respecting quotes
func parseCommandLine(input string) []string {
	var args []string
	var current strings.Builder
	var inQuotes bool
	var quoteChar rune

	input = strings.TrimSpace(input)
	runes := []rune(input)

	for i := 0; i < len(runes); i++ {
		char := runes[i]

		switch {
		case char == '"' || char == '\'':
			if !inQuotes {
				// Start of quoted section
				inQuotes = true
				quoteChar = char
			} else if char == quoteChar {
				// End of quoted section
				inQuotes = false
				quoteChar = 0
			} else {
				// Different quote type inside quotes
				current.WriteRune(char)
			}

		case char == '\\' && inQuotes && i+1 < len(runes):
			// Handle escape sequences inside quotes
			next := runes[i+1]
			switch next {
			case 'n':
				current.WriteRune('\n')
			case 't':
				current.WriteRune('\t')
			case 'r':
				current.WriteRune('\r')
			case '\\':
				current.WriteRune('\\')
			case '"', '\'':
				current.WriteRune(next)
			default:
				// Unknown escape, keep both characters
				current.WriteRune(char)
				current.WriteRune(next)
			}
			i++ // Skip the next character as we've processed it

		case unicode.IsSpace(char) && !inQuotes:
			// Space outside quotes - end current argument
			if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}
			// Skip multiple spaces
			for i+1 < len(runes) && unicode.IsSpace(runes[i+1]) {
				i++
			}

		default:
			// Regular character
			current.WriteRune(char)
		}
	}

	// Add the last argument if there is one
	if current.Len() > 0 {
		args = append(args, current.String())
	}

	return args
}

func HandleExternalCommand(cmd *Command) error {
	if cmd.command == "" {
		return nil
	}

	c := exec.Command(cmd.command, cmd.args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin

	if err := c.Run(); err != nil {
		return fmt.Errorf("traSH: %s: %v", cmd.command, err)
	}
	return nil
}

func HandleCommand(cmd *Command) error {
	if cmd.command == "" {
		return nil
	}

	switch cmd.command {
	case "cd":
		return HandleCD(cmd)
	case "exit", "quit":
		return nil
	case "help", "?":
		printHelp()
		return nil
	case "!ai":
		return HandleAI(cmd)
	case "!explain":
		return HandleExplain(cmd)
	default:
		return HandleExternalCommand(cmd)
	}
}
