package gosplit

import (
	"fmt"
	"strings"
)

type Game struct {
	gameName string
	runCategory string
	attempts string
}

type Split struct {
	SegmentName string
	SplitTime string
	SegmentTime string
	BestSegment string
}

var cmdTreeRoot *Command
var Splits []Split

func ParseCommands(cmd string) (string, error) {
	if cmdTreeRoot == nil {
		cmdTreeRoot = initCommandTree()
	}
	cmdSplit := strings.Split(cmd, " ")
	command, exists := cmdTreeRoot.subcommands[cmdSplit[0]]
	if !exists {
		return "", fmt.Errorf(
			"Unknown command: %s\nType \"help\" for a list of commands.",
			cmd,
		)
	}
	return execute(command, cmdSplit[1:])
}

func execute(cmd *Command, args []string) (string, error) {
	// base case
	if len(args) == 0 {
		if cmd.handler != nil {
			return cmd.handler([]string{})
		}
	}

	// help command needs a special case
	if cmd.name == cmdTreeRoot.subcommands["help"].name {
		return cmd.handler(args)
	}

	next, ok := cmd.subcommands[args[0]]
	if !ok {
		if cmd.handler != nil {
			return cmd.handler(args)
		}
		return "", fmt.Errorf("Unknown command: %s", args[0])
	}

	return execute(next, args[1:])
}
