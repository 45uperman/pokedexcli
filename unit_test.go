package main

import (
	"bytes"
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
		input    string
		expected []string
	}{
		{
			input: "help",
			expected: []string{
				"Welcome to the Pokedex!",
				"Usage:",
				"",
				"exit: Exit the Pokedex",
				"help: Displays a help message",
				"map: Prints the next page of areas",
			},
		},
		{
			input:    "map",
			expected: firstPageResults,
		},
		{
			input:    "mapb",
			expected: []string{"map is the command to open the map or move to the next page if the map is open, mapb is for moving to the previous page"},
		},
	}

	for _, c := range cases {
		testConfig := &config{}
		testConfig.RuntimeCache = pokecache.NewCache(30 * time.Second)
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w
		supportedCommands[c.input].callback(testConfig)
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
