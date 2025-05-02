package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/45uperman/pokedexcli/internal/farfetched"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

type config struct {
	Next      string
	Previous  string
	MapOpened bool
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
	runtimeConfig := &config{}
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
		err = commandStruct.callback(runtimeConfig)
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

func commandExit(cfg *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	isRunning = false
	return nil
}

func commandHelp(cfg *config) error {
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")
	for _, command := range supportedCommands {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}

func commandMap(cfg *config) error {
	var err error
	var path string

	if !cfg.MapOpened {
		path = farfetched.PokeURL + farfetched.LocationArea
	} else if cfg.Next == "" {
		fmt.Println("you're on the last page")
	} else {
		path = cfg.Next
	}

	data, err := farfetched.PokeGet(path)
	if err != nil {
		return err
	}

	mapPage, err := farfetched.BuildPage(data)
	if err != nil {
		return err
	}

	for _, result := range mapPage.Results {
		fmt.Println(result.Name)
	}

	if mapPage.Next == nil {
		cfg.Next = ""
	} else {
		cfg.Next = *mapPage.Next
	}
	if mapPage.Previous == nil {
		cfg.Previous = ""
	} else {
		cfg.Previous = *mapPage.Previous
	}
	cfg.MapOpened = true

	return nil
}

func commandMapB(cfg *config) error {
	var err error

	if !cfg.MapOpened {
		fmt.Println("map is the command to open the map or move to the next page if the map is open, mapb is for moving to the previous page")
		return nil
	} else if cfg.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}

	data, err := farfetched.PokeGet(cfg.Previous)
	if err != nil {
		return err
	}

	mapPage, err := farfetched.BuildPage(data)
	if err != nil {
		return err
	}

	for _, result := range mapPage.Results {
		fmt.Println(result.Name)
	}

	cfg.Next = *mapPage.Next
	if mapPage.Previous == nil {
		cfg.Previous = ""
	} else {
		cfg.Previous = *mapPage.Previous
	}

	return nil
}
