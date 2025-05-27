package command

import (
	"errors"
	"fmt"
	"os"
	"strings"
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

func ParseCommand(command string) *Command {
	commandList := strings.Fields(command)
	cmd := Command{
		command: commandList[0],
		args:    commandList[1:],
	}

	var filteredArgs []string
	for _, arg := range cmd.args {
		if len(arg) > 0 && arg[0] == '-' {
			cmd.opts = append(cmd.opts, []rune(arg[1:])...)
		} else {
			filteredArgs = append(filteredArgs, arg)
		}
	}
	cmd.args = filteredArgs

	return &cmd
}

func HandleCommand(cmd *Command) error {
	switch cmd.command {
	case "cd":
		return HandleCD(cmd)
	case "exit":
		os.Exit(1)
		return nil
	default:
		return errors.New("invalid command entered! -- use help for a list of available commands")
	}
}
