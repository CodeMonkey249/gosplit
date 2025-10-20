package gosplit

import (
	"bufio"
	"os"
	"encoding/json"
)

func ParseConfig(filepath string) []map[string]string {
	file, err := os.Open(filepath)
	Check(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var games []map[string]string

	for scanner.Scan() {
		line := scanner.Text()
		var data map[string]string

		err := json.Unmarshal([]byte(line), &data)
		Check(err)

		games = append(games, data)
	}

	return games
}
