package main

import (
	"bufio"
	"fmt"
	"gosplit/src"
	"os"
	"strings"
)

func main() {

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

