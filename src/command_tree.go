package gosplit

import (
	"fmt"
	//"slices"
	"strings"
)

type Command struct {
	name string
	cmd []string
	handler func([]string)(string, error)
	parent *Command
	subcommands map[string]*Command
	documentation string
}

func addSubCommand(cmd *Command, subCmd *Command) {
	for _, alias := range subCmd.cmd {
		cmd.subcommands[alias] = subCmd
	}
}

func buildDocumentation(root *Command) {

}

func initCommandTree() (*Command) {
	root := &Command {
		name: "ROOT",
		cmd: []string{"ROOT"},
		parent: nil,
		subcommands: map[string]*Command{},
		documentation: "Root command node, for organizational reasons only.",
	}

	newGame := &Command{
		name: "newGame",
		cmd: []string{"new_game", "ng"},
		handler: newGameHandler,
		parent: root,
		subcommands: map[string]*Command{},
		documentation: 
`
Configure a new game for GoSplit to manage.

You are given a series of prompts regarding game info. Your responses are
recorded and used to configure the game.

The prompts include the type of value expected (string, number, etc.). If
you provided an incorrect type, this command may fail.

Usage:	new_game | ng
`,
	}
	addSubCommand(root, newGame)

	list := &Command{
		name: "list",
		cmd: []string{"list", "l"},
		handler: listHandler,
		documentation: strings.Join([]string{
			"List all games.",
		}, "\n"),
		parent: root,
	}
	addSubCommand(root, list)

	_select := &Command{
		name: "select",
		cmd: []string{"select", "s"},
		handler: selectHandler,
		parent: root,
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
		parent: root,
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
		parent: root,
		documentation:
`
Perform game operations.

With no arguments, display the selected game.
`,
		subcommands: map[string]*Command{},
	}
	addSubCommand(root, game)

	split := &Command{
		name: "splits",
		cmd: []string{"splits", "sp"},
		handler: splitHandler,
		parent: root,
		documentation: 
`
Perform operations on splits.
`,
		subcommands: map[string]*Command{},
	}
	addSubCommand(root, split)

	addSplit := &Command{
		name: "addSplit",
		cmd: []string{"add"},
		handler: addSplitHandler,
		parent: split,
		documentation:
`
Add a new split

The first argument is the name to give the split.

The second argument is an optional index at which to add the split.

Usage:  game split add [SPLIT_NAME]
	game split add [SPLIT_NAME] [INDEX]
`,
	}
	addSubCommand(split, addSplit)

	removeSplit := &Command{
		name: "removeSplit",
		cmd: []string{"remove", "rm"},
		handler: removeSplitHandler,
		parent: split,
		documentation:
`
Remove a split.

Without any arguments, the current split will be removed.

Provide the index of the split you want to remove. The top split is 1, the second split is 2, etc.
Alternatively, you can provide the split name.

Usage:  split remove | split rm
	split remove [INDEX] | split remove [SPLIT_NAME]
`,
	}
	addSubCommand(split, removeSplit)

	start := &Command{
		name: "start",
		cmd: []string{"start", "go", "run"},
		handler: startHandler,
		parent: game,
		documentation: strings.Join([]string{
			"Start the timer.",
		}, "\n"),
	}
	addSubCommand(game, start)

	pause := &Command{
		name: "pause",
		cmd: []string{"pause", "p"},
		handler: pauseHandler,
		parent: game,
		documentation: strings.Join([]string{
			"Pause the timer.",
		}, "\n"),
	}
	addSubCommand(game, pause)

	stop := &Command{
		name: "stop",
		cmd: []string{"stop"},
		handler: stopHandler,
		parent: game,
		documentation: strings.Join([]string{
			"Stop the timer.",
		}, "\n"),
	}
	addSubCommand(game, stop)

	help := &Command{
		name: "help",
		cmd: []string{"help", "h"},
		handler: usageHandler,
		parent: root,
		subcommands: map[string]*Command{},
		documentation: strings.Join([]string{
			"Print list of commands.",
		}, "\n"),
	}
	addSubCommand(root, help)
	for _, cmd := range root.subcommands {
		addSubCommand(help, cmd)
	}

	edges := map[*Command][]*Command {
		root: {newGame, list, _select, quit, game, help, split},
		split: {removeSplit},
		game: {start, pause, stop},
		help: {newGame, list, _select, quit, game},
	}

	DFS(root, edges, visit)

	return root
}

func visit(cmd *Command) {
	if cmd.name == "ROOT" || cmd.parent.name == "ROOT" {
		return
	}

	msg := ""
	if !strings.Contains(
		cmd.parent.documentation,
		"List of " + cmd.parent.name + " subcommands",
	) {
		msg = fmt.Sprintf("\n\nList of %s subcommands:\n", cmd.parent.name)
	}

	absoluteCmd := buildCmdFromRoot(cmd)
	doc := strings.TrimSpace(cmd.documentation)
	firstLine := strings.Split(doc, "\n")[0]
	msg += fmt.Sprintf("\033[1m%s\033[0m -- %s\n", absoluteCmd, firstLine)

	cmd.parent.documentation += msg
}

func buildCmdFromRoot(cmd *Command) string {
	if cmd.parent.name == "ROOT" {
		return cmd.cmd[0]
	}
	return buildCmdFromRoot(cmd.parent) + " " + cmd.cmd[0]
}

func DFS(cmd *Command, edges map[*Command][]*Command, visitCb func(*Command)) {
	visited := map[string]bool{}

	if cmd == nil {
		return
	}
	visited[cmd.name] = true
	visitCb(cmd)

	for _, subCmd := range edges[cmd] {
		if visited[subCmd.name] {
			continue
		} else if subCmd.name == "help" {
			continue
		}
		DFS(subCmd, edges, visitCb)
	}
}
