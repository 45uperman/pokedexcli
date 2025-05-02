package main

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/45uperman/pokedexcli/internal/farfetched"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

var supportedCommands map[string]cliCommand
var isRunning bool
var mapPage farfetched.PokePage

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
		"map": {
			name:        "map",
			description: "Prints the next page of areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Prints the previous page of areas",
			callback:    commandMapB,
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
			break
		}
		input := cleanInput(scanner.Text())
		command := input[0]
		commandStruct, ok := supportedCommands[command]
		if !ok {
			fmt.Println("Invalid command!")
			continue
		}
		err = commandStruct.callback()
		if scanner.Err() != nil {
			fmt.Println(scanner.Err())
			isRunning = false
		}
		if err != nil {
			fmt.Println(err)
			isRunning = false
		}
	}
	if scanner.Err() != nil || err != nil {
		os.Exit(1)
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

func commandMap() error {
	var err error
	var path string

	if reflect.DeepEqual(mapPage, farfetched.PokePage{}) {
		path = farfetched.PokeURL + farfetched.LocationArea
	} else if mapPage.Next == nil {
		fmt.Println("you're on the last page")
		return nil
	} else {
		path = *mapPage.Next
	}

	err = farfetched.Farfetch(&mapPage, path)
	if err != nil {
		return err
	}

	for _, result := range mapPage.Results {
		fmt.Println(result.Name)
	}

	return nil
}

func commandMapB() error {
	var err error

	if reflect.DeepEqual(mapPage, farfetched.PokePage{}) {
		fmt.Println("map is the command to open the map for the first time and go forward through the pages, mapb goes back a page")
		return nil
	} else if mapPage.Previous == nil {
		fmt.Println("you're on the first page")
		return nil
	}

	err = farfetched.Farfetch(&mapPage, *mapPage.Previous)
	if err != nil {
		return err
	}

	for _, result := range mapPage.Results {
		fmt.Println(result.Name)
	}

	return nil
}
