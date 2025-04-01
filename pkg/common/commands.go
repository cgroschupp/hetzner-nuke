package common

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"
)

var commands []*cli.Command

// Commander --
type Commander interface {
	Execute(c context.Context)
}

// RegisterCommand --
func RegisterCommand(command *cli.Command) {
	logrus.Debugln("Registering", command.Name, "command...")
	commands = append(commands, command)
}

// GetCommands --
func GetCommands() []*cli.Command {
	return commands
}
