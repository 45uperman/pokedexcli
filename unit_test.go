package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/45uperman/pokedexcli/internal/farfetched"
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
	var err error
	var res *http.Response
	tries := 0
	for tries < 3 {
		res, err = http.Get("https://pokeapi.co/api/v2/location-area")
		if err == nil && res != nil {
			break
		}
		if err != nil {
			fmt.Printf("Encountered error while fetching data from PokeApi: %s\n", err)
		}
		tries++
	}
	if err != nil || res == nil {
		t.Errorf("Failed to fetch data from PokeApi after 3 tries\n")
		t.FailNow()
	}
	defer res.Body.Close()

	var firstMapPage farfetched.PokePage
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&firstMapPage)
	if err != nil {
		t.Errorf("Failed to stream PokeAPI response data to PokePage with error: %s", err)
		t.FailNow()
	}

	var mapResults []string
	for _, result := range firstMapPage.Results {
		fmt.Println(result)
		mapResults = append(mapResults, result.Name)
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
			expected: mapResults,
		},
	}

	for _, c := range cases {
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w
		supportedCommands[c.input].callback()
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
