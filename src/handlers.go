package gosplit

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"slices"
	"strconv"
	"text/tabwriter"
	"time"

	"strings"

	"github.com/MarinX/keylogger"
	"github.com/jedib0t/go-pretty/v6/table"
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

func padRight(s string, width int) string {
	return fmt.Sprintf("%-*s", width, s)
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func hideCursor() { fmt.Print("\033[?25l") }
func showCursor() { fmt.Print("\033[?25h") }

// moveCursorUp moves the cursor up by n lines
func moveCursorUp(x int, y int) {
	fmt.Printf("\033[%d;%dH", y, x)
}

func startHandler(args []string) (string, error) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	style := table.StyleDefault
	style.Box = table.BoxStyle{
		PaddingLeft:  " ",
		PaddingRight: " ",
	}
	style.Options = table.Options{
		DrawBorder:      false,
		SeparateColumns: true,
		SeparateHeader:  false,
		SeparateRows:    false,
	}
	t.SetStyle(style)

	ticker := time.NewTicker(1 * time.Millisecond)
	defer ticker.Stop()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	split := make(chan int, 1)
	go nextSplitHandler(split)

	start := time.Now()
	hideCursor()
	defer showCursor()

	clearScreen()
	moveCursorUp(1, 1)
	
	modSegments := []Segment{}
	for i := range Splits {
		modSegments = append(modSegments, Segment{})
		modSegments[i].SegmentName = Splits[i].SegmentName
		if Splits[i].SplitTime == normalizeTimeString("0") {
			modSegments[i].SplitTime = "-"
		} else {
			modSegments[i].SplitTime = Splits[i].SplitTime
		}
	}
	currSeg := 0

	for {
		select {
		case <-ticker.C:
			elapsed := time.Since(start).Round(time.Millisecond)

			t.ResetHeaders()
			t.ResetRows()
			t.AppendHeader(table.Row{
				padRight("Segment", 30),
				padRight("Time", 30),
			})

			for i, seg := range modSegments {
				if i == currSeg {
					t.AppendRows([]table.Row{
						{padRight(seg.SegmentName, 30), padRight(normalizeTimeString(elapsed.String()), 30)},
					})
				} else {
					t.AppendRows([]table.Row{
						{padRight(seg.SegmentName, 30), padRight(seg.SplitTime, 30)},
					})
				}
			}

			t.Render()
			moveCursorUp(1, 1)

		case <-stop:
			showCursor()
			return "", nil

		case <-split:
			elapsed := time.Since(start).Round(time.Millisecond)
			modSegments[currSeg].SplitTime = normalizeTimeString(elapsed.String())
			currSeg++
			if currSeg >= len(modSegments) {
				clearScreen()
				showCursor()
				printSplits()
				fmt.Printf("Would you like to save these times? (y/n): ")
				var save string
				fmt.Scan(&save)
				switch save {
				case "y", "Y":
					saveSegments(modSegments)
					printSplits()
					return "Segments saved!\n", nil
				case "n", "N":
					fmt.Printf("Are you sure? Type n again to discard times: ")
					var confirm string
					fmt.Scan(&confirm)
					if confirm == "n" || confirm == "N" {
						return "Segments discarded", nil
					}
				}
			}
			go nextSplitHandler(split)
		}
	}
}

func nextSplitHandler(notify chan int) {
	var keyboard *keylogger.KeyLogger
	key_device_path := FindKeyboard()
	keyboard, err := keylogger.New(key_device_path)
	if err != nil {
		panic(err)
	}
	defer keyboard.Close()
	ListenForKeystroke(keyboard, notify, "SPACE")
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
