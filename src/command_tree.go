package gosplit

import (
	"strings"
)

type Command struct {
	name string
	cmd []string
	documentation string
	handler func([]string)(string, error)
	subcommands map[string]*Command
}

func addSubCommand(cmd *Command, subCmd *Command) {
	for _, alias := range subCmd.cmd {
		cmd.subcommands[alias] = subCmd
	}
}

func initCommandTree() (*Command) {
	root := &Command {
		name: "ROOT",
		cmd: []string{"ROOT"},
		documentation: "Root command node, for organizational reasons only.",
		subcommands: map[string]*Command{},
	}

	newGame := &Command{
		name: "newGame",
		cmd: []string{"new_game", "ng"},
		handler: newGameHandler,
		documentation: 
`
Configure a new game for GoSplit to manage.

You are given a series of prompts regarding game info. Your responses are
recorded and used to configure the game.

The prompts include the type of value expected (string, number, etc.). If
you provided an incorrect type, this command may fail.

Usage:	new_game | ng
`,
		subcommands: map[string]*Command{},
	}
	addSubCommand(root, newGame)

	list := &Command{
		name: "list",
		cmd: []string{"list", "l"},
		handler: listHandler,
		documentation: strings.Join([]string{
			"List all games.",
		}, "\n"),
	}
	addSubCommand(root, list)

	_select := &Command{
		name: "select",
		cmd: []string{"select", "s"},
		handler: selectHandler,
		documentation: 
`
Select a game to split.

The first argument is the name (or the starting letters) of the game.
If no argument is provided, a list of games GoSplit is configured for is
displayed, and you may search (i.e. fuzzy find) the game you want.

Usage:	select | select [GAME]
`,
	}
	addSubCommand(root, _select)

	quit := &Command{
		name: "quit",
		cmd: []string{"quit", "exit", "q"},
		handler: quitHandler,
		documentation:
`
Exit GoSplit.

The first argument is optional and is used as the exit code for the GoSplit process.

Usage: 	quit | q | exit
	quit [EXIT_CODE] | q [EXIT_CODE] | exit [EXIT_CODE]
`,
	}
	addSubCommand(root, quit)

	game := &Command{
		name: "game",
		cmd: []string{"game", "g"},
		handler: gameHandler,
		documentation: strings.Join([]string{
			"Show currently selected game.",
		}, "\n"),
		subcommands: map[string]*Command{},
	}
	addSubCommand(root, game)

	start := &Command{
		name: "start",
		cmd: []string{"start", "go", "run"},
		handler: startHandler,
		documentation: strings.Join([]string{
			"Start the timer.",
		}, "\n"),
	}
	addSubCommand(root, start)
	addSubCommand(game, start)

	pause := &Command{
		name: "pause",
		cmd: []string{"pause", "p"},
		handler: pauseHandler,
		documentation: strings.Join([]string{
			"Pause the timer.",
		}, "\n"),
	}
	addSubCommand(root, pause)
	addSubCommand(game, pause)

	stop := &Command{
		name: "stop",
		cmd: []string{"stop"},
		handler: stopHandler,
		documentation: strings.Join([]string{
			"Stop the timer.",
		}, "\n"),
	}
	addSubCommand(root, stop)
	addSubCommand(game, stop)

	split := &Command{
		name: "split",
		cmd: []string{"split", "sp"},
		handler: splitHandler,
		documentation: strings.Join([]string{
			"Save the current time to the current split.",
		}, "\n"),
	}
	addSubCommand(root, split)
	addSubCommand(game, split)

	help := &Command{
		name: "help",
		cmd: []string{"help", "h"},
		handler: usageHandler,
		documentation: strings.Join([]string{
			"Print list of commands.",
		}, "\n"),
		subcommands: map[string]*Command{},
	}
	for _, cmd := range root.subcommands {
		addSubCommand(help, cmd)
	}
	addSubCommand(root, help)

	return root
}

