package main

import (
	"bufio"
	"fmt"
	"gosplit/src"
	"os"
	"strings"

	//"github.com/MarinX/keylogger"
)

func main() {
	// TODO: Boiler plate code for setting up global hotkeys
	//var keyboard *keylogger.KeyLogger
	//key_device_path := gosplit.FindKeyboard()
	//keyboard, err := keylogger.New(key_device_path)
	//gosplit.Check(err)
	//defer keyboard.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("(gosplit) ")
		cmd, err := reader.ReadString('\n')
		gosplit.Check(err)
		cmd = strings.TrimSpace(cmd)
		ret, err := gosplit.ParseCommands(cmd)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println(ret)
		}
	}
}

