package command

import (
	"fmt"
	"strings"

	"github.com/mush1e/traSH/utils"
)

type Command struct {
	command string
	args    []string
	opts    []rune
}

func ParseCommand(command string) *Command {
	commandList := strings.Fields(command)
	cmd := Command{
		command: commandList[0],
		args:    commandList[1:],
	}

	for idx, arg := range cmd.args {
		if arg[0] == '-' {
			cmd.opts = append(cmd.opts, []rune(arg[1:])...)
			utils.Remove(cmd.args, idx)
		}
	}

	fmt.Printf("%+v\n", cmd)

	return &cmd
}
