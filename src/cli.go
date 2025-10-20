package gosplit

import (
	"encoding/json"
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strings"
)

type Game struct {
	gameName string
	runCategory string
	attempts string
}

var cmdTreeRoot *Command

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
			cmd.handler(args)
		}
		return "", fmt.Errorf("Unknown command: %s", args[0])
	}

	return execute(next, args[1:])
}

func usageHandler(args []string) (string, error) {
	var msg string
	if len(args) > 0 {
		newCmd := cmdTreeRoot.subcommands["help"]
		oldCmd := cmdTreeRoot.subcommands["help"]
		exists := false
		for i := range args {
			oldCmd = newCmd
			newCmd, exists = oldCmd.subcommands[args[i]]
			if exists {
				msg = newCmd.documentation
			} else {
				return "", fmt.Errorf("Unknown command: %s", strings.Join(args, " "))
			}
		}
	} else {
		msg = `
List of commands:

help, h -- Display this menu
list, l -- List all games
quit, q -- Quit GoSplit
select, s -- Select a game to split

Type "help" followed by the command for full documentation.`
	}

	msg = strings.TrimSpace(msg)
	return msg, nil
}

func listHandler(args []string) (string, error) {
	games := ParseConfig("config.jsonl") // TODO: parametrize me
	msg := ""
	for i := range games {
		msg += "\n"
		for key, value := range games[i] {
			msg += key + ": " + value + "\n"
		}
	}
	return msg, nil
}

func selectHandler(args []string) (string, error) {
	return "", nil
}

func quitHandler(args []string) (string, error) {
	//
	// TODO: add ability to take optional exit code
	//
	os.Exit(0)
	return "", nil  // Never happens
}

func newGameHandler(args []string) (string, error) {
	reader := bufio.NewReader(os.Stdin)

	typ := reflect.TypeOf(Game{})

	newGameMap := map[string]string{}

	for i := 0; i < typ.NumField(); i++ {
		fieldName := typ.Field(i).Name
		fmt.Printf("%s: ", fieldName)
		fieldValue, err := reader.ReadString('\n')
		Check(err)
		fieldValue = strings.TrimSpace(fieldValue)
		newGameMap[fieldName] = fieldValue
	}

	//
	// TODO: input validation should be added here prior to writing to json.
	//

	filepath := "config_test.jsonl"  // TODO: parametrize me
	writeOutToJson(newGameMap, filepath)

	return "", nil
}

func gameHandler(args []string) (string, error) {
	return "", nil
}

func startHandler(args []string) (string, error) {
	return "", nil
}

func pauseHandler(args []string) (string, error) {
	return "", nil
}

func stopHandler(args []string) (string, error) {
	return "", nil
}

func splitHandler(args []string) (string, error) {
	return "", nil
}

func writeOutToJson(data map[string]string, filepath string) (error) {
	file, err := os.OpenFile(
		filepath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return err
	}
	defer file.Close()

    encoder := json.NewEncoder(file)
	encoder.Encode(data)

	return nil
}

// Args we need
// -h/--help
//    usage
//
// list
//   splits
// newGame
// newSplit
// select
// game
//   start
//   pause
//   stop
//   split

