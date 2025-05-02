package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/45uperman/pokedexcli/internal/farfetched"
	"github.com/45uperman/pokedexcli/internal/pokecache"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input: "    diglet  DIG!   diglet DIG!  TRIO!   TRIO!!   TRIO!!!     ",
			expected: []string{
				"diglet", "dig!", "diglet", "dig!", "trio!", "trio!!", "trio!!!",
			},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf(
				"Output of unexpected length\nEXPECTED: %v\nACTUAL: %v",
				c.expected,
				actual,
			)
			t.FailNow()
		}

		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf(
					"Found incorrect element!\nEXPECTED: %v\nACTUAL: %v",
					c.expected,
					actual,
				)
				t.FailNow()
			}
		}
	}
}

func TestCommands(t *testing.T) {
	testCache := pokecache.NewCache(30 * time.Second)
	data, err := farfetched.PokeGet(farfetched.PokeURL+farfetched.LocationArea, &testCache)
	if err != nil {
		t.Errorf("PokeGet failed to fetch data with error: %s", err)
		t.FailNow()
	}

	firstMapPage, err := farfetched.BuildPage(data)
	if err != nil {
		t.Errorf("BuildPage failed to build page with error: %s", err)
		t.FailNow()
	}

	var firstPageResults []string
	for _, result := range firstMapPage.Results {
		firstPageResults = append(firstPageResults, result.Name)
	}

	cases := []struct {
		input    []string
		expected []string
	}{
		{
			input: []string{"help"},
			expected: []string{
				"Welcome to the Pokedex!",
				"Usage:",
				"",
				"exit: Exit the Pokedex",
				"help: Displays a help message",
				"map: Displays the next page of areas",
				"mapb: Displays the previous page of areas",
				"explore: Displays all the Pokemon in the provided area",
				"catch: Tries to catch the provided Pokemon",
				"inspect: Displays information about the provided Pokemon",
			},
		},
		{
			input:    []string{"map"},
			expected: firstPageResults,
		},
		{
			input:    []string{"mapb"},
			expected: []string{"map is the command to open the map or move to the next page if the map is open, mapb is for moving to the previous page"},
		},
		{
			input: []string{"explore", "pastoria-city-area"},
			expected: []string{
				"Exploring pastoria-city-area...",
				"Found Pokemon:",
				" - tentacool",
				" - tentacruel",
			},
		},
		{
			input: []string{"catch", "squirtle"},
			// squirtle squad
			expected: []string{"Throwing a Pokeball at squirtle..."},
		},
		{
			input: []string{"inspect", "squirtle"},
			// squirtle squad
			expected: []string{" - defense: 65", " - special-defense: 64"},
		},
		{
			input:    []string{"pokedex"},
			expected: []string{" - squirtle"},
			// squirtle squad
		},
	}

	testConfig := &config{}
	testConfig.RuntimeCache = pokecache.NewCache(30 * time.Second)

	fullURL := fmt.Sprintf("%s/pokemon/%s", farfetched.PokeURL, "squirtle")

	data, err = farfetched.PokeGet(fullURL, &testConfig.RuntimeCache)
	if err != nil {
		t.Error("couldn't get squirtle D:")
		t.FailNow()
	}

	squirtle, err := farfetched.BuildPokemon("squirtle", fullURL, data)
	if err != nil {
		t.Error("couldn't get squirtle D:")
		t.FailNow()
	}

	for _, c := range cases {
		testConfig := &config{}
		testConfig.RuntimeCache = pokecache.NewCache(30 * time.Second)
		testConfig.Pokedex = farfetched.NewPokedex()
		testConfig.Pokedex.Catch(squirtle)
		var args []string
		if len(c.input) == 1 {
			args = []string{""}
		} else {
			args = c.input[1:]
		}
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w
		supportedCommands[c.input[0]].callback(testConfig, args)
		w.Close()
		os.Stdout = oldStdout
		var buf bytes.Buffer
		io.Copy(&buf, r)
		actual := buf.String()
		for _, line := range c.expected {
			if !strings.Contains(actual, line) {
				t.Errorf(
					"Missing line {%s} in output\nEXPECTED: %v\nACTUAL: %v",
					line,
					c.expected,
					actual,
				)
				t.FailNow()
			}
		}
	}
}
