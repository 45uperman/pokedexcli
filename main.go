package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/45uperman/pokedexcli/internal/farfetched"
	"github.com/45uperman/pokedexcli/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config, []string) error
}

type config struct {
	Next         string
	Previous     string
	MapOpened    bool
	RuntimeCache pokecache.Cache
	Pokedex      farfetched.Pokedex
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
			description: "Displays the next page of areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous page of areas",
			callback:    commandMapB,
		},
		"explore": {
			name:        "explore",
			description: "Displays all the Pokemon in the provided area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Tries to catch the provided Pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Displays information about the provided Pokemon",
			callback:    commandInspect,
		},
	}
	isRunning = true
}

func main() {
	var err error
	runtimeConfig := &config{}
	runtimeConfig.RuntimeCache = pokecache.NewCache(5 * time.Second)
	runtimeConfig.Pokedex = farfetched.NewPokedex()
	scanner := bufio.NewScanner(os.Stdin)
	for isRunning {
		fmt.Print("Pokedex > ")
		ok := scanner.Scan()
		if !ok {
			fmt.Println(scanner.Err())
			isRunning = false
			continue
		}
		input := cleanInput(scanner.Text())
		if len(input) == 0 {
			continue
		}

		command := input[0]
		var args []string
		if len(input) == 1 {
			args = []string{""}
		} else {
			args = input[1:]
		}

		commandStruct, ok := supportedCommands[command]
		if !ok {
			fmt.Println("Invalid command!")
			continue
		}

		err = commandStruct.callback(runtimeConfig, args)
		if err != nil {
			fmt.Println(err)
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

func commandExit(cfg *config, params []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	isRunning = false
	return nil
}

func commandHelp(cfg *config, params []string) error {
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")
	for _, command := range supportedCommands {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}

func commandMap(cfg *config, params []string) error {
	var err error
	var url string

	if !cfg.MapOpened {
		url = farfetched.PokeURL + farfetched.LocationArea
	} else if cfg.Next == "" {
		fmt.Println("you're on the last page")
	} else {
		url = cfg.Next
	}

	data, err := farfetched.PokeGet(url, &cfg.RuntimeCache)
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

func commandMapB(cfg *config, params []string) error {
	if !cfg.MapOpened {
		fmt.Println("map is the command to open the map or move to the next page if the map is open, mapb is for moving to the previous page")
		return nil
	} else if cfg.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}

	data, err := farfetched.PokeGet(cfg.Previous, &cfg.RuntimeCache)
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

func commandExplore(cfg *config, params []string) error {
	if params[0] == "" {
		fmt.Println("exlpore requires an area to explore - try looking at the map to see where you can go!")
		return nil
	}
	baseURL := farfetched.PokeURL + farfetched.LocationArea

	data, err := farfetched.PokeGet(
		fmt.Sprintf("%s/%s", baseURL, params[0]),
		&cfg.RuntimeCache,
	)
	if err != nil {
		return err
	}

	if string(data) == "Not Found" {
		return fmt.Errorf("error getting area info: resource %s not found", params[0])
	}

	encounters, err := farfetched.ExploreArea(data)
	if err != nil {
		return err
	}

	prefix := " - "
	fmt.Printf("Exploring %s...\n", params[0])
	fmt.Println("Found Pokemon:")
	for _, encounter := range encounters {
		fmt.Println(prefix + encounter.Pokemon.Name)
	}

	return nil
}

func commandCatch(cfg *config, params []string) error {
	if params[0] == "" {
		fmt.Println("catch requires a pokemon to catch - try exploring an area to see what's out there!")
		return nil
	}
	fullURL := fmt.Sprintf("%s/pokemon/%s", farfetched.PokeURL, params[0])

	data, err := farfetched.PokeGet(fullURL, &cfg.RuntimeCache)
	if err != nil {
		return err
	}

	if string(data) == "Not Found" {
		return fmt.Errorf("error getting pokemon info: resource %s not found", params[0])
	}

	pkmn, err := farfetched.BuildPokemon(params[0], fullURL, data)
	if err != nil {
		return err
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pkmn.Name)
	dice := rand.New(rand.NewSource(time.Now().UnixNano()))
	if dice.Intn(pkmn.BaseExperience) > 40 {
		fmt.Printf("%s escaped!\n", pkmn.Name)
	} else {
		cfg.Pokedex.Catch(pkmn)
		fmt.Printf("%s was caught!\n", pkmn.Name)
	}

	return nil
}

func commandInspect(cfg *config, params []string) error {
	if params[0] == "" {
		fmt.Println("inspect requires a pokemon to inspect - try checking your pokedex to see what pokemon you have!")
		return nil
	}
	pkmn := cfg.Pokedex.Pokemon[params[0]]
	indent := " - "

	fmt.Printf("Name: %s\n", pkmn.Name)
	fmt.Printf("Height: %d\n", pkmn.Height)
	fmt.Printf("Weight: %d\n", pkmn.Weight)
	fmt.Printf("ID: %d\n", pkmn.ID)

	fmt.Println("Stats:")
	for _, stat := range pkmn.Stats {
		fmt.Printf("%s%s: %d\n", indent, stat.Stat.Name, stat.BaseStat)
	}

	fmt.Println("Types:")
	for _, t := range pkmn.Types {
		fmt.Printf("%s%s\n", indent, t.Type.Name)
	}

	return nil
}
