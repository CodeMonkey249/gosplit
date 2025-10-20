package gosplit

import (
	"testing"
	//"regexp"
)

type cliTest struct {
	cmds []string
	expectation string
}

var cliTests = []cliTest {
	{
		cmds: []string{"help"},
		expectation: `List of commands:

help, h -- Display this menu
list, l -- List all games
quit, q -- Quit GoSplit
select, s -- Select a game to split

Type "help" followed by the command for full documentation.`,
	},
	{
		cmds: []string{"h"},
		expectation: `List of commands:

help, h -- Display this menu
list, l -- List all games
quit, q -- Quit GoSplit
select, s -- Select a game to split

Type "help" followed by the command for full documentation.`,
	},
	{
		cmds: []string{"help list"},
		expectation: `List all games.`,
	},
	{
		cmds: []string{"help select"},
		expectation: `Select a game to split.`,
	},
	{
		cmds: []string{"help quit"},
		expectation: `Exit GoSplit.
Usage: quit [EXIT_CODE] | exit [EXIT_CODE]`,
	},
	{
		cmds: []string{"help new_game"},
		expectation: `Create a new game listing.`,
	},
}

func TestBasicCLI(t *testing.T) {
	for _, test := range cliTests {
		output := ""
		for _, cmd := range test.cmds {
			ret, err := ParseCommands(cmd)
			output += ret
			if err != nil {
				t.Errorf("Error received: %s", err.Error())
			}
		}
		if output != test.expectation {
			t.Errorf(
				"Output doesn't match expectation.\nGot:\n%s\nExpected:\n%s",
				output,
				test.expectation,
			)
		}
	}
}
