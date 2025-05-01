package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

var supportedCommands map[string]cliCommand
var isRunning bool

func init() {
	supportedCommands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
	}
	isRunning = true
}

func main() {
	var err error

	scanner := bufio.NewScanner(os.Stdin)
	for isRunning {
		fmt.Print("Pokedex > ")
		ok := scanner.Scan()
		if !ok {
			if scanner.Err() != nil {
				err = scanner.Err()
			}
			break
		}
		input := cleanInput(scanner.Text())
		command := input[0]
		commandStruct, ok := supportedCommands[command]
		if !ok {
			fmt.Println("Invalid command!")
			continue
		}
		commandStruct.callback()
	}
	if err != nil {
		fmt.Println(err)
	}
}

func cleanInput(text string) (cleanWords []string) {
	trimmedText := strings.TrimSpace(text)
	words := strings.Split(trimmedText, " ")
	for i := range words {
		if strings.TrimSpace(words[i]) != "" {
			cleanWords = append(cleanWords, strings.ToLower(words[i]))
		}
	}
	return cleanWords
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	isRunning = false
	return nil
}

func commandHelp() error {
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")

	for _, command := range supportedCommands {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}
