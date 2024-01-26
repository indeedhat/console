package main

import (
	"bufio"
	"fmt"
	"os"
)

// cliPromptBool prompts the user for a boolean input from stdin
func cliPromptBool(prompt string) bool {
	for {
		fmt.Println(prompt + " [y/N]")

		data, err := readCliInput()
		if err != nil {
			continue
		}

		if len(data) == 0 {
			return false
		}

		if len(data) != 1 {
			continue
		}

		return data[0] == 'y' || data[0] == 'Y'
	}
}

// readInput reads a line of input from stdin
// it removes the trailing \n
func readCliInput() ([]byte, error) {
	reader := bufio.NewReader(os.Stdin)
	data, err := reader.ReadBytes('\n')

	if err != nil {
		return nil, err
	}

	return data[:len(data)-1], nil
}
