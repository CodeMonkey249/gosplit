package gosplit

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"slices"
	"strconv"
	"text/tabwriter"

	"strings"
)

func writeOutToJson(data map[string]string, filepath string, mode int) error {
	file, err := os.OpenFile(
		filepath,
		mode,
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
		helpCommand := cmdTreeRoot.subcommands["help"]
		helpSubcommands := helpCommand.subcommands
		helpSubcommandsKeys := make([]string, 0, len(helpSubcommands))
		for k := range helpSubcommands {
			helpSubcommandsKeys = append(helpSubcommandsKeys, k)
		}

		slices.Sort(helpSubcommandsKeys)

		cmdsCompleted := []*Command{}

		msg = "List of commands:\n\n"
		for _, key := range helpSubcommandsKeys {
			cmd := helpSubcommands[key]
			if slices.Contains(cmdsCompleted, cmd) {
				continue
			}

			doc := cmd.documentation
			doc = strings.TrimSpace(doc)
			firstLine := strings.Split(doc, "\n")[0]
			msg += strings.Join(cmd.cmd, ",")
			msg += " -- "
			msg += firstLine + "\n"

			cmdsCompleted = append(cmdsCompleted, cmd)
		}
		msg += "\nType \"help\" followed by the command for full documentation."
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
	return "", nil // Never happens
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
	filepath := "config_test.jsonl" // TODO: parametrize me
	mode := os.O_APPEND|os.O_CREATE|os.O_WRONLY
	writeOutToJson(newGameMap, filepath, mode)

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
	printSplits()
	return "", nil
}

func addSplitHandler(args []string) (string, error) {
	if len(args) <= 0 {
		return "", fmt.Errorf("You must provide a segment name.")
	}

	newSplit := Segment{
		SegmentName: args[0],
		SplitTime: normalizeTimeString("0"),
		SegmentTime: normalizeTimeString("0"),
		BestSegment: normalizeTimeString("99:99:99.999"),
	}
	splitIndex := len(Splits)
	for _, arg := range args {
		if strings.HasPrefix(arg, "at") {
			i := strings.Split(arg, "=")[1]
			err := fmt.Errorf("")
			splitIndex, err = strconv.Atoi(i)  // User defined index is 1-indexed
			splitIndex -= 1
			if err != nil {
				return "", fmt.Errorf("Error parsing user defined index: %v", err)
			}
		}
	}
	Splits = slices.Insert(Splits, splitIndex, newSplit)

	for _, arg := range args {
		if strings.HasPrefix(arg, "SplitTime") {
			splitTime := strings.Split(arg, "=")[1]
			splitTime = normalizeTimeString(splitTime)
			err := setSplitTime(splitTime, splitIndex)
			if err != nil {
				return "", err
			}
		}
		if strings.HasPrefix(arg, "SegmentTime") {
			segmentTime := strings.Split(arg, "=")[1]
			segmentTime = normalizeTimeString(segmentTime)
			err := setSegmentTime(segmentTime, splitIndex)
			if err != nil {
				return "", err
			}
		}
		if strings.HasPrefix(arg, "BestSegment") {
			bestSegment := strings.Split(arg, "=")[1]
			bestSegment = normalizeTimeString(bestSegment)
			setBestSegment(bestSegment, splitIndex)
		}
	}

	printSplits()
	return "", nil
}

func printSplits() {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 5, ' ', 0)

	fmt.Fprintf(w, "\nSegment\tSplit Time\tSegment Time\tBest Segment\n")
	for _, sp := range Splits {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", sp.SegmentName, sp.SplitTime, sp.SegmentTime, sp.BestSegment)
	}
	w.Flush()
}

func removeSplitHandler(args []string) (string, error) {
	return "", nil
}
